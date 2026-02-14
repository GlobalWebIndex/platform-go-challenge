package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gwi/platform-go-challenge/favourites/application"
	"github.com/gwi/platform-go-challenge/favourites/domain"
)

type OpenAPIServer struct {
	UC           application.FavouritesUseCase
	EventsUC     application.FavouritesUseCase
	EventsSecret string
	Signer       Signer
	Verifier     Verifier
	TokenExpiry  time.Duration
}

var _ ServerInterface = (*OpenAPIServer)(nil)

func intPtr(i int) *int       { return &i }
func strPtr(s string) *string { return &s }

const apiVersion = "1"
const apiName = "favourites"

func (s *OpenAPIServer) GetVersion(w http.ResponseWriter, r *http.Request) {
	resp := VersionResponse{Version: strPtr(apiVersion), Api: strPtr(apiName)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *OpenAPIServer) GetHealth(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	status := Healthy
	statusCode := http.StatusOK

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	done := make(chan bool, 1)
	go func() {
		// Perform a lightweight query to verify system is working
		_, _ = s.UC.ListFavouritesPaginated("health-check", 1, 0)
		done <- true
	}()

	select {
	case <-done:
	case <-ctx.Done():
		// Health check timed out or was cancelled
		status = Unhealthy
		statusCode = http.StatusServiceUnavailable
		slog.Warn("health check failed", "error", "timeout")
	}

	resp := HealthResponse{
		Status:    &status,
		Timestamp: &now,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *OpenAPIServer) authorizeUser(w http.ResponseWriter, r *http.Request, pathUserID string) bool {
	return AuthorizeUser(w, r, s.Verifier, pathUserID)
}

func (s *OpenAPIServer) IssueToken(w http.ResponseWriter, r *http.Request) {
	var req IssueTokenJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorWithLog(w, r, http.StatusBadRequest, ErrCodeBadRequest, "invalid JSON body", err)
		return
	}
	if msg := validateVar(req.UserId, "required"); msg != "" {
		WriteError(w, http.StatusBadRequest, ErrCodeBadRequest, "user_id: "+msg)
		return
	}
	expiry := s.TokenExpiry
	if expiry <= 0 {
		expiry = DefaultTokenExpiry
	}
	now := time.Now()
	exp := now.Add(expiry)
	claims := &jwt.RegisteredClaims{
		Subject:   req.UserId,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(exp),
	}
	tokenStr, err := s.Signer.Sign(claims)
	if err != nil {
		WriteErrorWithLog(w, r, http.StatusInternalServerError, ErrCodeInternal, "failed to issue token", err)
		return
	}
	sec := int(expiry.Seconds())
	resp := IssueTokenResponse{Token: &tokenStr, ExpiresIn: &sec}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *OpenAPIServer) AssetDeleted(w http.ResponseWriter, r *http.Request) {
	if s.EventsSecret != "" && strings.TrimSpace(r.Header.Get("X-Internal-Secret")) != s.EventsSecret {
		WriteError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "invalid or missing X-Internal-Secret")
		return
	}
	var req AssetDeletedJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid JSON body")
		return
	}
	if msg := validateVar(req.AssetId, "required"); msg != "" {
		WriteError(w, http.StatusBadRequest, ErrCodeBadRequest, "asset_id: "+msg)
		return
	}
	removed := s.EventsUC.OnAssetDeleted(req.AssetId)
	resp := AssetDeletedResponse{Removed: &removed}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *OpenAPIServer) ListFavourites(w http.ResponseWriter, r *http.Request, userID UserID, params ListFavouritesParams) {
	if !s.authorizeUser(w, r, string(userID)) {
		return
	}
	limit, offset := DefaultListLimit, 0
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
		if limit > MaxListLimit {
			limit = MaxListLimit
		}
	}
	if params.Offset != nil && *params.Offset >= 0 {
		offset = *params.Offset
	}
	items, total := s.UC.ListFavouritesPaginated(string(userID), limit, offset)
	listItems := make([]FavouriteListItem, 0, len(items))
	for _, a := range items {
		listItems = append(listItems, FavouriteListItem{
			Id:            strPtr(a.ID),
			Type:          strPtr(a.Type),
			Description:   strPtr(a.Description),
			SourceAssetId: strPtr(a.SourceAssetID),
			CreatedAt:     &a.CreatedAt,
			UpdatedAt:     &a.UpdatedAt,
		})
	}
	resp := ListFavouritesResponse{
		Items:  &listItems,
		Total:  &total,
		Limit:  &limit,
		Offset: &offset,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *OpenAPIServer) AddFavourite(w http.ResponseWriter, r *http.Request, userID UserID) {
	if !s.authorizeUser(w, r, string(userID)) {
		return
	}
	var req AddFavouriteJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorWithLog(w, r, http.StatusBadRequest, ErrCodeBadRequest, "invalid JSON body", err)
		return
	}
	if msg := validateVar(req.Type, "required,oneof=chart insight audience"); msg != "" {
		errMsg := "type: " + msg
		if strings.Contains(msg, "oneof") {
			errMsg = "type must be one of: chart, insight, audience"
		}
		WriteError(w, http.StatusBadRequest, ErrCodeBadRequest, errMsg)
		return
	}
	if err := validateAddFavouriteRequest(req); err != "" {
		WriteError(w, http.StatusBadRequest, ErrCodeBadRequest, err)
		return
	}
	asset := apiAddRequestToDomain(req)
	added, err := s.UC.AddFavourite(string(userID), asset)
	if err != nil {
		if errors.Is(err, application.ErrAssetNotFound) {
			WriteError(w, http.StatusBadRequest, ErrCodeBadRequest, "asset not found in catalog: cannot add favourite for unknown source_asset_id")
			return
		}
		WriteErrorWithLog(w, r, http.StatusInternalServerError, ErrCodeInternal, "failed to add favourite", err)
		return
	}
	resp := domainToFavouriteResponse(added)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *OpenAPIServer) GetFavourite(w http.ResponseWriter, r *http.Request, userID UserID, assetID AssetID) {
	if !s.authorizeUser(w, r, string(userID)) {
		return
	}
	asset, ok := s.UC.GetFavourite(string(userID), string(assetID))
	if !ok {
		WriteError(w, http.StatusNotFound, ErrCodeNotFound, "asset not found")
		return
	}
	resp := domainToFavouriteResponse(*asset)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *OpenAPIServer) RemoveFavourite(w http.ResponseWriter, r *http.Request, userID UserID, assetID AssetID) {
	if !s.authorizeUser(w, r, string(userID)) {
		return
	}
	if !s.UC.RemoveFavourite(string(userID), string(assetID)) {
		WriteError(w, http.StatusNotFound, ErrCodeNotFound, "asset not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *OpenAPIServer) UpdateFavouriteDescription(w http.ResponseWriter, r *http.Request, userID UserID, assetID AssetID) {
	if !s.authorizeUser(w, r, string(userID)) {
		return
	}
	var req UpdateFavouriteDescriptionJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorWithLog(w, r, http.StatusBadRequest, ErrCodeBadRequest, "invalid JSON body", err)
		return
	}
	desc := ""
	if req.Description != nil {
		desc = *req.Description
		if len(desc) > 500 {
			WriteError(w, http.StatusBadRequest, ErrCodeBadRequest, "description must be 500 characters or less")
			return
		}
	}
	updated, ok := s.UC.UpdateFavouriteDescription(string(userID), string(assetID), desc)
	if !ok {
		WriteError(w, http.StatusNotFound, ErrCodeNotFound, "asset not found")
		return
	}
	resp := domainToFavouriteResponse(*updated)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func validateAddFavouriteRequest(req AddFavouriteJSONRequestBody) string {
	if req.Description != nil && len(*req.Description) > 500 {
		return "description must be 500 characters or less"
	}

	addByRef := req.SourceAssetId != nil && *req.SourceAssetId != ""

	switch req.Type {
	case domain.AssetTypeChart:
		if !addByRef && req.Chart == nil {
			return "type \"chart\" requires a \"chart\" object (or set source_asset_id to add by reference)"
		}
		if req.Chart != nil {
			if req.Chart.Title != nil && strings.TrimSpace(*req.Chart.Title) == "" {
				return "chart title cannot be empty"
			}
			if req.Chart.XAxisTitle != nil && strings.TrimSpace(*req.Chart.XAxisTitle) == "" {
				return "chart x_axis_title cannot be empty"
			}
			if req.Chart.YAxisTitle != nil && strings.TrimSpace(*req.Chart.YAxisTitle) == "" {
				return "chart y_axis_title cannot be empty"
			}
			if req.Chart.Title != nil && len(*req.Chart.Title) > 200 {
				return "chart title must be 200 characters or less"
			}
		}

	case domain.AssetTypeInsight:
		if !addByRef && req.Insight == nil {
			return "type \"insight\" requires an \"insight\" object (or set source_asset_id to add by reference)"
		}
		if req.Insight != nil {
			if req.Insight.Text != nil && strings.TrimSpace(*req.Insight.Text) == "" {
				return "insight text cannot be empty"
			}
			if req.Insight.Text != nil && len(*req.Insight.Text) > 2000 {
				return "insight text must be 2000 characters or less"
			}
		}

	case domain.AssetTypeAudience:
		if !addByRef && req.Audience == nil {
			return "type \"audience\" requires an \"audience\" object (or set source_asset_id to add by reference)"
		}
		if req.Audience != nil {
			if req.Audience.HoursSocialMediaDaily != nil && *req.Audience.HoursSocialMediaDaily < 0 {
				return "audience hours_social_media_daily cannot be negative"
			}
			if req.Audience.HoursSocialMediaDaily != nil && *req.Audience.HoursSocialMediaDaily > 24 {
				return "audience hours_social_media_daily cannot exceed 24 hours"
			}
			if req.Audience.PurchasesLastMonth != nil && *req.Audience.PurchasesLastMonth < 0 {
				return "audience purchases_last_month cannot be negative"
			}
		}

	default:
		return "invalid asset type: must be chart, insight, or audience"
	}
	return ""
}

func apiAddRequestToDomain(req AddFavouriteJSONRequestBody) domain.Asset {
	a := domain.Asset{Type: req.Type}
	if req.Description != nil {
		a.Description = *req.Description
	}
	if req.SourceAssetId != nil {
		a.SourceAssetID = *req.SourceAssetId
	}
	if req.Chart != nil {
		a.Chart = &domain.ChartData{}
		if req.Chart.Title != nil {
			a.Chart.Title = *req.Chart.Title
		}
		if req.Chart.XAxisTitle != nil {
			a.Chart.XAxisTitle = *req.Chart.XAxisTitle
		}
		if req.Chart.YAxisTitle != nil {
			a.Chart.YAxisTitle = *req.Chart.YAxisTitle
		}
		if req.Chart.Data != nil {
			a.Chart.Data = req.Chart.Data
		}
	}
	if req.Insight != nil && req.Insight.Text != nil {
		a.Insight = &domain.InsightData{Text: *req.Insight.Text}
	}
	if req.Audience != nil {
		a.Audience = &domain.AudienceData{}
		if req.Audience.Gender != nil {
			a.Audience.Gender = *req.Audience.Gender
		}
		if req.Audience.BirthCountry != nil {
			a.Audience.BirthCountry = *req.Audience.BirthCountry
		}
		if req.Audience.AgeGroups != nil {
			a.Audience.AgeGroups = *req.Audience.AgeGroups
		}
		if req.Audience.HoursSocialMediaDaily != nil {
			a.Audience.HoursSocialMediaDaily = float64(*req.Audience.HoursSocialMediaDaily)
		}
		if req.Audience.PurchasesLastMonth != nil {
			a.Audience.PurchasesLastMonth = *req.Audience.PurchasesLastMonth
		}
	}
	return a
}

func domainToFavouriteResponse(a domain.Asset) FavouriteResponse {
	resp := FavouriteResponse{
		Id:            strPtr(a.ID),
		Type:          strPtr(a.Type),
		Description:   strPtr(a.Description),
		SourceAssetId: strPtr(a.SourceAssetID),
		CreatedAt:     &a.CreatedAt,
		UpdatedAt:     &a.UpdatedAt,
	}
	if a.Chart != nil {
		resp.Chart = &ChartData{
			Title:      strPtr(a.Chart.Title),
			XAxisTitle: strPtr(a.Chart.XAxisTitle),
			YAxisTitle: strPtr(a.Chart.YAxisTitle),
			Data:       a.Chart.Data,
		}
	}
	if a.Insight != nil {
		resp.Insight = &InsightData{Text: strPtr(a.Insight.Text)}
	}
	if a.Audience != nil {
		h := float32(a.Audience.HoursSocialMediaDaily)
		resp.Audience = &AudienceData{
			Gender:                strPtr(a.Audience.Gender),
			BirthCountry:          strPtr(a.Audience.BirthCountry),
			AgeGroups:             strPtr(a.Audience.AgeGroups),
			HoursSocialMediaDaily: &h,
			PurchasesLastMonth:    &a.Audience.PurchasesLastMonth,
		}
	}
	return resp
}

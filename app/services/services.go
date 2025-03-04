package services

import (
	"github.com/biter777/countries"

	"app/api/schema"
	"app/models"
	"app/models/constants"
	"app/repositories"

	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)


type GService interface {
	FindAssets(filter schema.AssetFilter) []schema.AssetDetailsDto
	FindAsset(assetID string) (schema.AssetDetailsDto, error)
	UpdateAsset(assetID string, description string) error
	FindUser(userID string) (models.User, error)
	FindUserFavorites(userID string) []models.UserFavorite
	CreateFavoriteAsset(assetID string, userID string) (string, error)
	DeleteUserFavorite(userFavoriteID string) error
}

type EntityService struct {
	Repo repositories.GRepository
}

func NewGService(repo repositories.GRepository) GService {
	return &EntityService{Repo: repo}
}

func (s *EntityService) FindAssets(filter schema.AssetFilter) []schema.AssetDetailsDto {
	assets := s.Repo.FindAssets(filter)
	if len(assets) == 0 {
		return make([]schema.AssetDetailsDto, 0)
	}
	response := make([]schema.AssetDetailsDto, 0, len(assets))
	for _, a := range assets {
		dto, err := getAssetDetailsDto(a)
		if err != nil {
			fmt.Printf("[!] failed to get asset details for asset %s due to error {%s}, skipping...\n", a, err.Error())
			continue
		} else {
			response = append(response, dto)
		}
	}
	return response
}

func (s *EntityService) FindAsset(assetID string) (schema.AssetDetailsDto, error) {
	assetVo, err := s.Repo.FindAsset(assetID)
	if err != nil {
		return schema.AssetDetailsDto{}, err
	}
	if assetVo.ID == "" {
		return schema.AssetDetailsDto{}, errors.New("asset not found")
	}
	return getAssetDetailsDto(assetVo)
}

func (s *EntityService) UpdateAsset(assetID string, description string) error {
	request := models.AssetUpdateRequest{AssetID: assetID, Description: description}
	return s.Repo.UpdateAsset(request)
}

func (s *EntityService) FindUser(userID string) (models.User, error) {
	return s.Repo.FindUser(userID)
}

func (s *EntityService) FindUserFavorites(userID string) []models.UserFavorite {
	return s.Repo.FindUserFavorites(userID)
}

func (s *EntityService) CreateFavoriteAsset(assetID string, userID string) (string, error) {
	user, err := s.FindUser(userID)
	if err != nil {
		fmt.Printf("user %s not found\n", userID)
		return "", schema.NewApiError(http.StatusNotFound, errors.Join(err, errors.New("user not found")))
	}
	userFavorites := s.Repo.FindUserFavorites(userID)

	for _, uf := range userFavorites {
		if uf.AssetID == assetID {
			errorMessage := fmt.Sprintf("asset %s is already a favorite for user %s", assetID, userID)
			return "", schema.NewApiError(http.StatusConflict, errors.New(errorMessage))
		}
	}
	asset, err := s.FindAsset(assetID)
	if err != nil {
		fmt.Printf("asset %s not found", assetID)
		return "", schema.NewApiError(http.StatusNotFound, errors.New("asset not found"))
	}
	return s.Repo.CreateUserFavorite(user.ID.Hex(), asset.Asset.ID), nil
}

func (s *EntityService) DeleteUserFavorite(userFavoriteID string) error {
	return s.Repo.DeleteUserFavorite(userFavoriteID)
}

func getAssetDetailsDto(vo models.AssetVO) (schema.AssetDetailsDto, error) {
	switch vo.Type {
	case constants.ASSET_TYPE_CHART:
		return schema.AssetDetailsDto{
			Asset:        schema.AssetDto{ID: vo.ID, Description: vo.Description},
			ChartDetails: &schema.ChartAssetDto{Title: vo.Title, AxesTitles: vo.AxesTitles, PlotData: vo.PlotData},
		}, nil
	case constants.ASSET_TYPE_INSIGHT:
		return schema.AssetDetailsDto{
			Asset:        schema.AssetDto{ID: vo.ID, Description: vo.Description},
			InsightDetails: &schema.InsightAssetDto{Text: vo.Text},
		}, nil
		case constants.ASSET_TYPE_AUDIENCE:
			var characteristics []string = make([]string, 0, len(vo.Characteristics))
			for _, c := range vo.Characteristics {
				characteristicsStr, err := getCharacteristicStr(c.Key, c.Value)
				if err != nil {
					return schema.AssetDetailsDto{}, err
				}
				characteristics = append(characteristics, characteristicsStr)
			}
			return schema.AssetDetailsDto{
				Asset:        schema.AssetDto{ID: vo.ID, Description: vo.Description},
				AudienceDetails: &schema.AudienceAssetDto{Characteristics: characteristics},
			}, nil
		default: return schema.AssetDetailsDto{}, errors.New(fmt.Sprintf("unable to map AssetVO to AssetDetailsDto for asset %s", vo.ID))
	}
}

func getCharacteristicStr(characteristicID int, characteristicValue int) (string, error) {
	var sb strings.Builder
	if characteristicID == constants.CHARACTERISTIC_AGE_GROUP {
		ageGroup := getAgeGroupText(characteristicValue)
		sb.WriteString(schema.CHARACTERISTIC_AGE_GROUP)
		sb.WriteString(": ")
		sb.WriteString(ageGroup)
	} else if characteristicID == constants.CHARACTERISTIC_BIRTH_COUNTRY {
		country := countries.ByNumeric(characteristicValue).Info().Name
		sb.WriteString(schema.CHARACTERISTIC_BIRTH_COUNTRY)
		sb.WriteString(": ")
		sb.WriteString(country)
	} else if characteristicID == constants.CHARACTERISTIC_GENDER {
		gender := getGenderText(characteristicValue)
		sb.WriteString(schema.CHARACTERISTIC_GENDER)
		sb.WriteString(": ")
		sb.WriteString(gender)
	} else if characteristicID == constants.CHARACTERISTIC_PURCHASES_LAST_MONTH {
		sb.WriteString(schema.CHARACTERISTIC_PURCHASES_LAST_MONTH)
		sb.WriteString(": ")
		sb.WriteString(strconv.Itoa(characteristicValue))
	} else if characteristicID == constants.CHARACTERISTIC_SOCIAL_MEDIA_DAILY_HOURS {
		sb.WriteString(schema.CHARACTERISTIC_SOCIAL_MEDIA_DAILY_HOURS)
		sb.WriteString(": ")
		sb.WriteString(strconv.Itoa(characteristicValue))
	} else {
		return "", errors.New("unhandled characteristic id")
	}
	return sb.String(), nil
} 

func getGenderText(genderID int) string {
	if genderID == 1 {
		return schema.GENDER_FEMALE
	}
	return schema.GENDER_MALE
}

func getAgeGroupText(ageGroupID int) string {
	switch ageGroupID {
	case constants.AGE_GROUP_17_19:
		return "17 to 19"
	case constants.AGE_GROUP_20_29:
		return "20 to 29"
	case constants.AGE_GROUP_30_39:
		return "30 to 39"
	case constants.AGE_GROUP_40_49:
		return "40 to 49"
	case constants.AGE_GROUP_50_59:
		return "50 to 59"
	case constants.AGE_GROUP_60_69:
		return "60 to 69"
	case constants.AGE_GROUP_70_PLUS:
		return "70+"
	default:
		return "UNKNOWN"
	}
}

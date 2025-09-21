package gwiexercise

import (
	"encoding/json"
	"testing"
)

func TestAddAssetRequest_ValidateBasic(t *testing.T) {
	tests := []struct {
		name     string
		request  AddAssetRequest
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "valid request",
			request: AddAssetRequest{Type: "chart", Description: "Test Description", Payload: json.RawMessage(`{}`)},
			wantErr: false,
		},
		{
			name:     "missing type",
			request:  AddAssetRequest{Description: "Test Description", Payload: json.RawMessage(`{}`)},
			wantErr:  true,
			errorMsg: "type is required",
		},
		{
			name:     "missing description",
			request:  AddAssetRequest{Type: "chart", Payload: json.RawMessage(`{}`)},
			wantErr:  true,
			errorMsg: "description is required",
		},
		{
			name:     "missing payload",
			request:  AddAssetRequest{Type: "chart", Description: "Test Description"},
			wantErr:  true,
			errorMsg: "payload is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.ValidateBasic()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBasic() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errorMsg {
				t.Errorf("ValidateBasic() error = %v, expected error = %v", err, tt.errorMsg)
			}
		})
	}
}

func TestDecodeAssetFromAddRequest(t *testing.T) {
	tests := []struct {
		name       string
		req        AddAssetRequest
		wantType   AssetType
		wantErr    bool
		errorMsg   string
		verifyFunc func(Asset) bool
	}{
		{
			name: "valid chart asset",
			req: AddAssetRequest{
				ID:          "1",
				Type:        "chart",
				Description: "Chart Description",
				Payload:     json.RawMessage(`{"title":"ChartTitle", "xAxisTitle":"XTitle", "yAxisTitle":"YTitle", "data":[1.0, 2.0, 3.0]}`),
			},
			wantType: "chart",
			wantErr:  false,
			verifyFunc: func(a Asset) bool {
				chart, ok := a.(*ChartAsset)
				return ok && chart.Title == "ChartTitle" && len(chart.Data) == 3
			},
		},
		{
			name: "valid insight asset",
			req: AddAssetRequest{
				ID:          "2",
				Type:        "insight",
				Description: "Insight Description",
				Payload:     json.RawMessage(`{"text":"InsightText"}`),
			},
			wantType: "insight",
			wantErr:  false,
			verifyFunc: func(a Asset) bool {
				insight, ok := a.(*InsightAsset)
				return ok && insight.Text == "InsightText"
			},
		},
		{
			name: "valid audience asset",
			req: AddAssetRequest{
				ID:          "3",
				Type:        "audience",
				Description: "Audience Description",
				Payload:     json.RawMessage(`{"gender":"Female","birthCountry":"US","ageGroup":"25-34","socialHoursPerDay":"2-3","purchasesLastMonth":5}`),
			},
			wantType: "audience",
			wantErr:  false,
			verifyFunc: func(a Asset) bool {
				audience, ok := a.(*AudienceAsset)
				return ok && audience.Gender == "Female" && audience.PurchasesLastMonth == 5
			},
		},
		{
			name: "invalid asset type",
			req: AddAssetRequest{
				ID:          "4",
				Type:        "unknown",
				Description: "Invalid Description",
				Payload:     json.RawMessage(`{}`),
			},
			wantErr:  true,
			errorMsg: "unsupported asset type: unknown",
		},
		{
			name: "invalid chart payload",
			req: AddAssetRequest{
				ID:          "5",
				Type:        "chart",
				Description: "Chart Description",
				Payload:     json.RawMessage(`{"invalidField":"value"}`),
			},
			wantErr:  true,
			errorMsg: "invalid chart payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset, err := DecodeAssetFromAddRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeAssetFromAddRequest() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !containsError(err.Error(), tt.errorMsg) {
					t.Errorf("DecodeAssetFromAddRequest() error = %v, expected error to contain = %v", err, tt.errorMsg)
				}
				return
			}
			if asset.GetType() != tt.wantType {
				t.Errorf("DecodeAssetFromAddRequest() type = %v, want = %v", asset.GetType(), tt.wantType)
			}
			if tt.verifyFunc != nil && !tt.verifyFunc(asset) {
				t.Errorf("DecodeAssetFromAddRequest() verification failed")
			}
		})
	}
}

func TestDecodeAssetFromJSON(t *testing.T) {
	tests := []struct {
		name       string
		jsonData   []byte
		wantType   AssetType
		wantErr    bool
		errorMsg   string
		verifyFunc func(Asset) bool
	}{
		{
			name: "valid chart JSON",
			jsonData: json.RawMessage(`{
				"id": "1",
				"type": "chart",
				"description": "Chart Description",
				"title": "Chart Title",
				"xAxisTitle": "X Axis",
				"yAxisTitle": "Y Axis",
				"data": [1.0, 2.0, 3.0]
			}`),
			wantType: "chart",
			wantErr:  false,
			verifyFunc: func(a Asset) bool {
				chart, ok := a.(*ChartAsset)
				return ok && chart.Title == "Chart Title" && len(chart.Data) == 3
			},
		},
		{
			name: "unsupported asset type",
			jsonData: json.RawMessage(`{
				"id": "2",
				"type": "unsupported",
				"description": "Invalid Asset"
			}`),
			wantErr:  true,
			errorMsg: "unsupported asset type: unsupported",
		},
		{
			name:     "invalid JSON data",
			jsonData: json.RawMessage(`{invalid JSON`),
			wantErr:  true,
			errorMsg: "invalid asset json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset, err := DecodeAssetFromJSON(tt.jsonData)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeAssetFromJSON() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !containsError(err.Error(), tt.errorMsg) {
					t.Errorf("DecodeAssetFromJSON() error = %v, expected error to contain = %v", err, tt.errorMsg)
				}
				return
			}
			if asset.GetType() != tt.wantType {
				t.Errorf("DecodeAssetFromJSON() type = %v, want = %v", asset.GetType(), tt.wantType)
			}
			if tt.verifyFunc != nil && !tt.verifyFunc(asset) {
				t.Errorf("DecodeAssetFromJSON() verification failed")
			}
		})
	}
}

func containsError(actual, expected string) bool {
	return len(actual) >= len(expected) && actual[:len(expected)] == expected
}

package favorites

import (
	"encoding/json"
	"testing"
)

func TestInsight_Validate(t *testing.T) {
	tests := []struct {
		name    string
		insight Insight
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid insight",
			insight: Insight{
				BaseAsset: BaseAsset{
					ID:   "insight-1",
					Name: "Market Trend",
					Type: AssetTypeInsight,
				},
				Content: "Market grew by 5%",
			},
			wantErr: false,
		},
		{
			name: "missing content",
			insight: Insight{
				BaseAsset: BaseAsset{
					ID:   "insight-1",
					Name: "Empty Insight",
					Type: AssetTypeInsight,
				},
			},
			wantErr: true,
			errMsg:  "validation failed: content is required for Insight",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.insight.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: Insight.Validate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("%s: Insight.Validate() error = %v, wantErrMsg %v", tt.name, err, tt.errMsg)
			}
		})
	}
}

func TestInsight_JSON(t *testing.T) {
	jsonData := `{
		"id": "ins-123",
		"type": "insight",
		"name": "Insight Name",
		"content": "Secret sauce"
	}`

	var ins Insight
	err := json.Unmarshal([]byte(jsonData), &ins)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if ins.Content != "Secret sauce" {
		t.Errorf("Expected content 'Secret sauce', got '%s'", ins.Content)
	}
}

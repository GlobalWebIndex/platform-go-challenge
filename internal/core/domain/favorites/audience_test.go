package favorites

import (
	"encoding/json"
	"testing"
)

func TestAudience_Validate(t *testing.T) {
	tests := []struct {
		name    string
		aud     Audience
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid audience",
			aud: Audience{
				BaseAsset: BaseAsset{
					ID:   "aud-1",
					Name: "Gen Z Europe",
					Type: AssetTypeAudience,
				},
				Rules: AudienceRules{
					Country: "UK",
					AgeMin:  18,
					AgeMax:  25,
				},
			},
			wantErr: false,
		},
		{
			name: "missing country",
			aud: Audience{
				BaseAsset: BaseAsset{
					ID:   "aud-1",
					Name: "Invalid",
					Type: AssetTypeAudience,
				},
				Rules: AudienceRules{
					AgeMin: 18,
				},
			},
			wantErr: true,
			errMsg:  "validation failed: country rule is required for Audience",
		},
		{
			name: "negative age",
			aud: Audience{
				BaseAsset: BaseAsset{
					ID:   "aud-1",
					Name: "Invalid",
					Type: AssetTypeAudience,
				},
				Rules: AudienceRules{
					Country: "UK",
					AgeMin:  -1,
				},
			},
			wantErr: true,
			errMsg:  "validation failed: age limits cannot be negative",
		},
		{
			name: "min > max age",
			aud: Audience{
				BaseAsset: BaseAsset{
					ID:   "aud-1",
					Name: "Invalid",
					Type: AssetTypeAudience,
				},
				Rules: AudienceRules{
					Country: "UK",
					AgeMin:  30,
					AgeMax:  20,
				},
			},
			wantErr: true,
			errMsg:  "validation failed: age_min cannot be greater than age_max",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.aud.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: Audience.Validate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("%s: Audience.Validate() error = %v, wantErrMsg %v", tt.name, err, tt.errMsg)
			}
		})
	}
}

func TestAudience_JSON(t *testing.T) {
	jsonData := `{
		"id": "aud-123",
		"type": "audience",
		"name": "Target Group",
		"rules": {
			"country": "US",
			"age_min": 21
		}
	}`

	var aud Audience
	err := json.Unmarshal([]byte(jsonData), &aud)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if aud.ID != "aud-123" {
		t.Errorf("Expected ID 'aud-123', got '%s'", aud.ID)
	}
	if aud.Rules.Country != "US" {
		t.Errorf("Expected country 'US', got '%s'", aud.Rules.Country)
	}
	if aud.Rules.AgeMin != 21 {
		t.Errorf("Expected age_min 21, got %d", aud.Rules.AgeMin)
	}
}

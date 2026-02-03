package favorites

import (
	"encoding/json"
	"testing"
)

func TestChart_Validate(t *testing.T) {
	tests := []struct {
		name    string
		chart   Chart
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid chart",
			chart: Chart{
				BaseAsset: BaseAsset{
					ID:   "chart-1",
					Name: "Revenue",
					Type: AssetTypeChart,
				},
				XAxis: "Month",
				YAxis: "USD",
			},
			wantErr: false,
		},
		{
			name: "missing both axes",
			chart: Chart{
				BaseAsset: BaseAsset{
					ID:   "chart-1",
					Name: "Empty Chart",
					Type: AssetTypeChart,
				},
			},
			wantErr: true,
			errMsg:  "validation failed: chart requires at least one axis to be defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.chart.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: Chart.Validate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("%s: Chart.Validate() error = %v, wantErrMsg %v", tt.name, err, tt.errMsg)
			}
		})
	}
}

func TestChart_JSON(t *testing.T) {
	jsonData := `{
		"id": "chart-99",
		"type": "chart",
		"name": "Bar Chart",
		"x_axis": "Category"
	}`

	var chart Chart
	err := json.Unmarshal([]byte(jsonData), &chart)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if chart.XAxis != "Category" {
		t.Errorf("Expected XAxis 'Category', got '%s'", chart.XAxis)
	}
	if chart.YAxis != "" {
		t.Errorf("Expected YAxis empty, got '%s'", chart.YAxis)
	}
}

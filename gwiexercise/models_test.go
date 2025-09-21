package gwiexercise

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBaseAsset(t *testing.T) {
	base := &BaseAsset{
		ID:          "123",
		Type:        "Chart",
		Description: "Test Description",
		CreatedAt:   time.Now(),
	}

	t.Run("GetID", func(t *testing.T) {
		if id := base.GetID(); id != "123" {
			t.Errorf("expected ID to be '123', got '%s'", id)
		}
	})

	t.Run("GetType", func(t *testing.T) {
		if typ := base.GetType(); typ != "Chart" {
			t.Errorf("expected Type to be 'Chart', got '%s'", typ)
		}
	})

	t.Run("GetDescription", func(t *testing.T) {
		if desc := base.GetDescription(); desc != "Test Description" {
			t.Errorf("expected Description to be 'Test Description', got '%s'", desc)
		}
	})

	t.Run("SetDescription", func(t *testing.T) {
		base.SetDescription("Updated Description")
		if desc := base.GetDescription(); desc != "Updated Description" {
			t.Errorf("expected Description to be 'Updated Description', got '%s'", desc)
		}
	})
}

func TestChartAsset_MarshalJSON(t *testing.T) {
	chart := &ChartAsset{
		BaseAsset: BaseAsset{
			ID:          "001",
			Type:        "Chart",
			Description: "Chart description",
			CreatedAt:   time.Now(),
		},
		Title:      "Sample Chart",
		XAxisTitle: "X Axis",
		YAxisTitle: "Y Axis",
		Data:       []float64{1.0, 2.5, 3.9},
	}

	bytes, err := chart.MarshalJSON()
	if err != nil {
		t.Fatalf("unexpected error during MarshalJSON: %v", err)
	}

	var unmarshalledChart ChartAsset
	if err := json.Unmarshal(bytes, &unmarshalledChart); err != nil {
		t.Fatalf("unexpected error during UnmarshalJSON: %v", err)
	}

	if unmarshalledChart.Title != chart.Title {
		t.Errorf("expected Title to be '%s', got '%s'", chart.Title, unmarshalledChart.Title)
	}
}

func TestInsightAsset_MarshalJSON(t *testing.T) {
	insight := &InsightAsset{
		BaseAsset: BaseAsset{
			ID:          "002",
			Type:        "Insight",
			Description: "Insight description",
			CreatedAt:   time.Now(),
		},
		Text: "Insightful text",
	}

	bytes, err := insight.MarshalJSON()
	if err != nil {
		t.Fatalf("unexpected error during MarshalJSON: %v", err)
	}

	var unmarshalledInsight InsightAsset
	if err := json.Unmarshal(bytes, &unmarshalledInsight); err != nil {
		t.Fatalf("unexpected error during UnmarshalJSON: %v", err)
	}

	if unmarshalledInsight.Text != insight.Text {
		t.Errorf("expected Text to be '%s', got '%s'", insight.Text, unmarshalledInsight.Text)
	}
}

func TestAudienceAsset_MarshalJSON(t *testing.T) {
	audience := &AudienceAsset{
		BaseAsset: BaseAsset{
			ID:          "003",
			Type:        "Audience",
			Description: "Audience description",
			CreatedAt:   time.Now(),
		},
		Gender:             "Female",
		BirthCountry:       "USA",
		AgeGroup:           "25-34",
		SocialHoursPerDay:  "3-4",
		PurchasesLastMonth: 5,
	}

	bytes, err := audience.MarshalJSON()
	if err != nil {
		t.Fatalf("unexpected error during MarshalJSON: %v", err)
	}

	var unmarshalledAudience AudienceAsset
	if err := json.Unmarshal(bytes, &unmarshalledAudience); err != nil {
		t.Fatalf("unexpected error during UnmarshalJSON: %v", err)
	}

	if unmarshalledAudience.Gender != audience.Gender {
		t.Errorf("expected Gender to be '%s', got '%s'", audience.Gender, unmarshalledAudience.Gender)
	}
}

package gwiexercise

import (
	"testing"
	"time"
)

func TestInMemoryStore_Add_AssignsID_AndSetsCreatedAtIfZero(t *testing.T) {
	s := NewInMemoryStore()

	// CreatedAt should be set if zero and ID should be auto-assigned if empty
	chart := &ChartAsset{BaseAsset: BaseAsset{Type: AssetChart, Description: "desc"}, Title: "T"}
	added, err := s.Add("u1", chart)
	if err != nil {
		t.Fatalf("Add returned error: %v", err)
	}
	if added.GetID() == "" {
		t.Fatalf("expected ID to be assigned")
	}
	if chart.CreatedAt.IsZero() {
		t.Fatalf("expected CreatedAt to be set")
	}

	// If ID is supplied, it should be respected
	ins := &InsightAsset{BaseAsset: BaseAsset{ID: "given-id", Type: AssetInsight, Description: "d"}, Text: "hi"}
	added2, err := s.Add("u1", ins)
	if err != nil {
		t.Fatalf("Add returned error: %v", err)
	}
	if added2.GetID() != "given-id" {
		t.Fatalf("expected given id to be preserved, got %q", added2.GetID())
	}
}

func TestHelpers_base36_fmtID(t *testing.T) {
	// base36 known values
	cases := []struct {
		in  uint64
		out string
	}{
		{0, "0"},
		{1, "1"},
		{10, "a"},
		{35, "z"},
		{36, "10"},
		{71, "1z"},
		{1295, "zz"},
	}
	for _, c := range cases {
		if got := base36(c.in); got != c.out {
			t.Fatalf("base36(%d) = %q, want %q", c.in, got, c.out)
		}
	}

	// fmtID should concatenate two base36 values
	x := int64(36) // base36 -> "10"
	y := int64(71) // base36 -> "1z"
	if got := fmtID(x, y); got != "101z" {
		t.Fatalf("fmtID unexpected: %q", got)
	}
}

func Test_setCreatedAtIfZero_onlySetsWhenZero(t *testing.T) {
	s := NewInMemoryStore()
	// Add one with existing CreatedAt, ensure it remains
	ts := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	chart := &ChartAsset{BaseAsset: BaseAsset{ID: "keep", Type: AssetChart, Description: "d", CreatedAt: ts}, Title: "t"}
	if _, err := s.Add("u", chart); err != nil {
		t.Fatalf("add: %v", err)
	}
	if chart.CreatedAt != ts {
		t.Fatalf("CreatedAt should remain unchanged when already set")
	}
}

package gwiexercise

import (
	"path/filepath"
	"sort"
	"testing"
)

type storeFactory func(t *testing.T) Store

func storesUnderTest() []struct {
	name string
	new  storeFactory
} {
	return []struct {
		name string
		new  storeFactory
	}{
		{
			name: "inmem",
			new:  func(t *testing.T) Store { return NewInMemoryStore() },
		},
		{
			name: "file",
			new: func(t *testing.T) Store {
				tmp := t.TempDir()
				p := filepath.Join(tmp, "store.json")
				fs, err := NewFileStore(p)
				if err != nil {
					t.Fatalf("NewFileStore: %v", err)
				}
				return fs
			},
		},
	}
}

func TestStoreConformance(t *testing.T) {
	for _, tc := range storesUnderTest() {
		// Duplicate ID within same user should error; same ID across users allowed
		t.Run(tc.name+"/AddDuplicateAndIsolation", func(t *testing.T) {
			st := tc.new(t)

			// add one with explicit ID
			if _, err := st.Add("u1", &InsightAsset{BaseAsset: BaseAsset{ID: "same", Type: AssetInsight, Description: "d"}, Text: "x"}); err != nil {
				t.Fatalf("add: %v", err)
			}
			// duplicate for same user should fail
			if _, err := st.Add("u1", &InsightAsset{BaseAsset: BaseAsset{ID: "same", Type: AssetInsight, Description: "d2"}, Text: "y"}); err == nil || err.Error() != "asset with same id already exists for user" {
				t.Fatalf("expected duplicate error, got %v", err)
			}
			// same id for another user is allowed
			if _, err := st.Add("u2", &InsightAsset{BaseAsset: BaseAsset{ID: "same", Type: AssetInsight, Description: "d3"}, Text: "z"}); err != nil {
				t.Fatalf("unexpected error for different user: %v", err)
			}
		})

		// List returns all for a user and is empty for non-existent users; order irrelevant
		t.Run(tc.name+"/List", func(t *testing.T) {
			st := tc.new(t)
			// empty for unknown user
			if got := st.List("nouser"); len(got) != 0 {
				t.Fatalf("expected empty list for unknown user, got %d", len(got))
			}
			// add two assets
			st.Add("u", &AudienceAsset{BaseAsset: BaseAsset{ID: "b", Type: AssetAudience, Description: "db"}})
			st.Add("u", &AudienceAsset{BaseAsset: BaseAsset{ID: "a", Type: AssetAudience, Description: "da"}})
			got := st.List("u")
			if len(got) != 2 {
				t.Fatalf("expected 2 assets, got %d", len(got))
			}
			ids := []string{got[0].GetID(), got[1].GetID()}
			sort.Strings(ids)
			if !(ids[0] == "a" && ids[1] == "b") {
				t.Fatalf("unexpected ids: %v", ids)
			}
		})

		// Remove success and error cases (missing asset/user)
		t.Run(tc.name+"/Remove", func(t *testing.T) {
			st := tc.new(t)
			st.Add("u", &ChartAsset{BaseAsset: BaseAsset{ID: "x", Type: AssetChart, Description: "d"}, Title: "t"})
			if err := st.Remove("u", "x"); err != nil {
				t.Fatalf("remove existing: %v", err)
			}
			if err := st.Remove("u", "missing"); err == nil || err.Error() != "asset not found" {
				t.Fatalf("expected asset not found, got %v", err)
			}
			if err := st.Remove("missing", "x"); err == nil || err.Error() != "user not found" {
				t.Fatalf("expected user not found, got %v", err)
			}
		})

		// UpdateDescription success and error cases
		t.Run(tc.name+"/UpdateDescription", func(t *testing.T) {
			st := tc.new(t)
			st.Add("u", &InsightAsset{BaseAsset: BaseAsset{ID: "id1", Type: AssetInsight, Description: "old"}, Text: "t"})
			updated, err := st.UpdateDescription("u", "id1", "new")
			if err != nil {
				t.Fatalf("update: %v", err)
			}
			if updated.GetDescription() != "new" {
				t.Fatalf("description not updated: %v", updated.GetDescription())
			}
			if _, err := st.UpdateDescription("nouser", "id1", "x"); err == nil || err.Error() != "user not found" {
				t.Fatalf("expected user not found, got %v", err)
			}
			if _, err := st.UpdateDescription("u", "missing", "x"); err == nil || err.Error() != "asset not found" {
				t.Fatalf("expected asset not found, got %v", err)
			}
		})
	}
}

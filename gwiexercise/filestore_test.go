package gwiexercise

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// fixedTime used to make timestamps deterministic in tests
var fixedTime = time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

func withFixedNow() func() {
	prev := now
	now = func() time.Time { return fixedTime }
	return func() { now = prev }
}

func TestFileStore_NewFileStore_NonexistentFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "store.json")

	fs, err := NewFileStore(path)
	if err != nil {
		t.Fatalf("NewFileStore returned error: %v", err)
	}
	if got := fs.List("any"); len(got) != 0 {
		t.Fatalf("expected empty list for fresh store, got %d", len(got))
	}
}

func TestFileStore_LoadsExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "store.json")

	// Prepare on-disk structure
	a1 := &ChartAsset{BaseAsset: BaseAsset{ID: "a1", Type: AssetChart, Description: "desc1", CreatedAt: fixedTime}, Title: "t", XAxisTitle: "x", YAxisTitle: "y", Data: []float64{1, 2}}
	a2 := &InsightAsset{BaseAsset: BaseAsset{ID: "a2", Type: AssetInsight, Description: "desc2", CreatedAt: fixedTime}, Text: "hello"}

	encUsers := map[string]map[string]json.RawMessage{}
	for uid, amap := range map[string]map[string]Asset{"u1": {"a1": a1}, "u2": {"a2": a2}} {
		inner := map[string]json.RawMessage{}
		for aid, asset := range amap {
			b, err := json.Marshal(asset)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			inner[aid] = b
		}
		encUsers[uid] = inner
	}
	disk := struct {
		Users map[string]map[string]json.RawMessage `json:"users"`
	}{Users: encUsers}
	fileBytes, _ := json.Marshal(disk)
	if err := os.WriteFile(path, fileBytes, 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	fs, err := NewFileStore(path)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	gotU1 := fs.List("u1")
	if len(gotU1) != 1 || gotU1[0].GetID() != "a1" {
		t.Fatalf("expected one asset a1 for u1, got %v", gotU1)
	}
	gotU2 := fs.List("u2")
	if len(gotU2) != 1 || gotU2[0].GetID() != "a2" {
		t.Fatalf("expected one asset a2 for u2, got %v", gotU2)
	}
}

func TestFileStore_Add_AssignsID_And_CreatedAt_And_Persists(t *testing.T) {
	cleanup := withFixedNow()
	defer cleanup()

	dir := t.TempDir()
	// nested path to ensure directories are created
	path := filepath.Join(dir, "nested", "data", "store.json")

	fs, err := NewFileStore(path)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	// asset without ID and CreatedAt
	chart := &ChartAsset{BaseAsset: BaseAsset{Type: AssetChart, Description: "d0"}, Title: "t"}
	added, err := fs.Add("user", chart)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if added.GetID() == "" {
		t.Fatalf("expected ID to be assigned")
	}
	// check CreatedAt on concrete type
	if chart.CreatedAt != fixedTime {
		t.Fatalf("expected CreatedAt %v, got %v", fixedTime, chart.CreatedAt)
	}

	// file should exist and contain our user
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	var disk struct {
		Users map[string]map[string]json.RawMessage `json:"users"`
	}
	if err := json.Unmarshal(b, &disk); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := disk.Users["user"]; !ok {
		t.Fatalf("expected user key in file")
	}

	// Reload store and ensure asset is present
	fs2, err := NewFileStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	got := fs2.List("user")
	if len(got) != 1 {
		t.Fatalf("expected 1 asset after reload, got %d", len(got))
	}
}

func TestFileStore_UpdateDescription_Persists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "store.json")
	fs, _ := NewFileStore(path)

	fs.Add("u", &ChartAsset{BaseAsset: BaseAsset{ID: "c1", Type: AssetChart, Description: "old"}, Title: "t"})

	updated, err := fs.UpdateDescription("u", "c1", "new-desc")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.GetDescription() != "new-desc" {
		t.Fatalf("description not updated: %v", updated.GetDescription())
	}

	// reload and verify
	fs2, err := NewFileStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	got := fs2.List("u")
	if len(got) != 1 || got[0].GetDescription() != "new-desc" {
		t.Fatalf("expected new-desc after reload, got %v", got)
	}
}

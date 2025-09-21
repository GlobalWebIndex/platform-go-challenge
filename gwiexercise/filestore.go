package gwiexercise

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileStore persists favourites to a single JSON file on disk.
// The file format is: { "users": { "<userID>": {"<assetID>": <asset-json>, ...}, ... } }
type FileStore struct {
	mu    sync.RWMutex
	path  string
	users map[string]map[string]Asset
}

func NewFileStore(path string) (*FileStore, error) {
	fs := &FileStore{path: path, users: make(map[string]map[string]Asset)}
	if err := fs.load(); err != nil {
		return nil, err
	}
	return fs, nil
}

func (s *FileStore) List(userID string) []Asset {
	s.mu.RLock()
	defer s.mu.RUnlock()
	umap := s.users[userID]
	if umap == nil {
		return []Asset{}
	}
	res := make([]Asset, 0, len(umap))
	for _, a := range umap {
		res = append(res, a)
	}
	return res
}

func (s *FileStore) ensureUser(userID string) map[string]Asset {
	m, ok := s.users[userID]
	if !ok {
		m = make(map[string]Asset)
		s.users[userID] = m
	}
	return m
}

func (s *FileStore) Add(userID string, a Asset) (Asset, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a.GetID() == "" {
		setID(a, newID())
	}
	setCreatedAtIfZero(a, now())
	umap := s.ensureUser(userID)
	if _, exists := umap[a.GetID()]; exists {
		return nil, errors.New("asset with same id already exists for user")
	}
	umap[a.GetID()] = a
	if err := s.saveLocked(); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *FileStore) Remove(userID, assetID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	umap := s.users[userID]
	if umap == nil {
		return errors.New("user not found")
	}
	if _, ok := umap[assetID]; !ok {
		return errors.New("asset not found")
	}
	delete(umap, assetID)
	return s.saveLocked()
}

func (s *FileStore) UpdateDescription(userID, assetID, description string) (Asset, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	umap := s.users[userID]
	if umap == nil {
		return nil, errors.New("user not found")
	}
	a, ok := umap[assetID]
	if !ok {
		return nil, errors.New("asset not found")
	}
	a.SetDescription(description)
	umap[assetID] = a
	if err := s.saveLocked(); err != nil {
		return nil, err
	}
	return a, nil
}

type onDisk struct {
	Users map[string]map[string]json.RawMessage `json:"users"`
}

func (s *FileStore) load() error {
	contents, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var disk onDisk
	if err := json.Unmarshal(contents, &disk); err != nil {
		return err
	}
	for uid, message := range disk.Users {
		umap := make(map[string]Asset, len(message))
		for aid, raw := range message {
			a, err := DecodeAssetFromJSON(raw)
			if err != nil {
				return err
			}
			umap[aid] = a
		}
		s.users[uid] = umap
	}
	return nil
}

func (s *FileStore) saveLocked() (err error) {
	d := onDisk{Users: make(map[string]map[string]json.RawMessage, len(s.users))}
	for uid, m := range s.users {
		mm := make(map[string]json.RawMessage, len(m))
		for aid, a := range m {
			b, err := json.Marshal(a)
			if err != nil {
				return err
			}
			mm[aid] = b
		}
		d.Users[uid] = mm
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(true)
	enc.SetIndent("", "  ")
	if e := enc.Encode(&d); e != nil {
		return e
	}
	return nil
}

// now is overridden in tests to make times deterministic.
var now = func() time.Time { return time.Now() }

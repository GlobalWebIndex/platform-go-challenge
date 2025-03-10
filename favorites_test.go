package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

const (
	invalidID = "1000000"
)

func TestFavoritesStorage(t *testing.T) {
	ownerA := "ownerA"
	ownerB := "ownerB"
	fs := newFavoritesStorage()

	t.Run("add", func(t *testing.T) {
		addedFavorite := Favorite{AssetKey: "keyA", Description: "descrA"}
		id := fs.add(ownerA, addedFavorite)

		if id == 0 {
			t.Errorf("add returned 0 id")
			t.FailNow()
		}

		// a 2nd add of the same favorite should be nop and return the same id
		if nid := fs.add(ownerA, addedFavorite); nid != id {
			t.Errorf("add added same favorite twice")
			t.FailNow()
		}
	})

	t.Run("getByID", func(t *testing.T) {
		addedFavorite := Favorite{AssetKey: "keyA", Description: "descrA"}
		id := fs.add(ownerA, addedFavorite)

		expectedFavorite := Favorite{
			ID:          id,
			AssetKey:    addedFavorite.AssetKey,
			Description: addedFavorite.Description,
		}

		getForA, present := fs.getByID(ownerA, id)
		if !present {
			t.Errorf("getByIndex did not return value")
			t.FailNow()
		}
		if getForA != expectedFavorite {
			t.Errorf("getByID mismatch: got %v expected %v", getForA, expectedFavorite)
			t.FailNow()
		}
		getForB, present := fs.getByID(ownerB, id)
		if present || getForB == expectedFavorite {
			t.Errorf("getByID leaks favorites between owners")
			t.FailNow()
		}
	})

	t.Run("getAll", func(t *testing.T) {
		ida := fs.add(ownerA, Favorite{AssetKey: "keyA", Description: "descrA"})
		idb := fs.add(ownerA, Favorite{AssetKey: "keyB", Description: "descrB"})
		idc := fs.add(ownerA, Favorite{AssetKey: "keyC", Description: "descrC"})
		fs.del(ownerA, idb)
		_ = fs.add(ownerB, Favorite{AssetKey: "keyD", Description: "descrD"})

		favs := fs.getAll(ownerA)
		if len(favs) != 2 {
			t.Errorf("getAll len mismatch. got %d expected %d", len(favs), 2)
			t.FailNow()
		}
		if favs[0].ID != ida {
			t.Errorf("favs[0] id mismatch. got %d expected %d", favs[0].ID, ida)
			t.FailNow()
		}
		if favs[1].ID != idc {
			t.Errorf("favs[1] id mismatch. got %d expected %d", favs[1].ID, idc)
			t.FailNow()
		}
	})

	t.Run("del", func(t *testing.T) {
		addedFavorite := Favorite{AssetKey: "keyA", Description: "descrA"}
		ida := fs.add(ownerA, addedFavorite)
		idb := fs.add(ownerB, addedFavorite)
		fs.del(ownerA, ida)
		if _, present := fs.getByID(ownerA, ida); present {
			t.Errorf("del failed, get still returns the value")
			t.FailNow()
		}
		if _, present := fs.getByID(ownerB, idb); !present {
			t.Errorf("del failed, deleted value for both owners")
			t.FailNow()
		}
	})

	t.Run("update", func(t *testing.T) {
		addedFavorite := Favorite{AssetKey: "keyA", Description: "descrBefore"}
		id := fs.add(ownerA, addedFavorite)
		updatedFavorite := Favorite{Description: "descrAfter"}
		fs.update(ownerA, id, updatedFavorite)

		expectedFavorite := Favorite{
			ID:          id,
			AssetKey:    addedFavorite.AssetKey,
			Description: updatedFavorite.Description,
		}

		getForA, present := fs.getByID(ownerA, id)
		if !present || getForA != expectedFavorite {
			t.Errorf("update failed. got %v expected %v", getForA, expectedFavorite)
			t.FailNow()
		}
	})
}

func TestFavoritesService(t *testing.T) {
	user := User{"test@gwifav.dev", "passwd"}

	t.Run("add", func(t *testing.T) {
		srv := newFavoritesService(newFavoritesStorage(), nil)

		t.Run("malformedInput", func(t *testing.T) {
			w := callAdd(srv, user, bytes.NewBufferString("hello"))
			if w.Code != http.StatusBadRequest {
				t.Errorf("add accepted malformed input")
			}
		})

		t.Run("wrongInput", func(t *testing.T) {
			addedFavorite := Favorite{Description: "descr"} // missing required AssetKey
			w := callAdd(srv, user, encodeJson(t, addedFavorite))
			if w.Code != http.StatusBadRequest {
				t.Errorf("add accepted wrong input")
			}
		})

		t.Run("correctInput", func(t *testing.T) {
			addedFavorite := Favorite{AssetKey: "key", Description: "descr"}
			w := callAdd(srv, user, encodeJson(t, addedFavorite))
			if w.Code != http.StatusCreated {
				t.Errorf("add rejected correct input")
			}
			if loc := w.Header().Get("Location"); loc == "" {
				t.Errorf("add did not return the favorites url")
			} else if !strings.HasPrefix(loc, "/favorites/") || !unicode.IsDigit(rune(loc[len(loc)-1])) {
				t.Errorf("add did return wrong favorites url: %s", loc)
			}
		})
	})

	t.Run("get", func(t *testing.T) {
		srv := newFavoritesService(newFavoritesStorage(), nil)
		addedFavorite := Favorite{AssetKey: "key", Description: "descr"}
		id := mustAddFavorite(t, srv, user, addedFavorite)

		t.Run("nonNumericID", func(t *testing.T) {
			w := callGet(srv, user, "hello")
			if w.Code != http.StatusNotFound {
				t.Errorf("get accepted non numeric id")
			}
		})

		t.Run("nonExistingFavorite", func(t *testing.T) {
			w := callGet(srv, user, invalidID)
			if w.Code != http.StatusNotFound {
				t.Errorf("get returned non existing favorite")
			}
		})

		t.Run("existingFavorite", func(t *testing.T) {
			returnedFavorite, ok := mustGetFavorite(t, srv, user, id)

			matchKey := addedFavorite.AssetKey == returnedFavorite.AssetKey
			matchDescription := addedFavorite.Description == returnedFavorite.Description
			if !ok || !matchKey || !matchDescription || returnedFavorite.ID == 0 {
				t.Errorf("get response error: got %v expected %v", returnedFavorite, addedFavorite)
			}
		})

	})

	t.Run("del", func(t *testing.T) {
		srv := newFavoritesService(newFavoritesStorage(), nil)
		addedFavorite := Favorite{AssetKey: "key", Description: "descr"}
		id := mustAddFavorite(t, srv, user, addedFavorite)

		t.Run("nonNumericID", func(t *testing.T) {
			w := callDelete(srv, user, "hello")
			if w.Code != http.StatusNotFound {
				t.Errorf("delete accepted non numeric id")
			}
		})

		t.Run("nonExistingFavorite", func(t *testing.T) {
			w := callDelete(srv, user, invalidID)
			if w.Code != http.StatusNotFound {
				t.Errorf("del accepted non existing favorite")
			}
		})

		t.Run("existingFavorite", func(t *testing.T) {
			w := callDelete(srv, user, id)
			if w.Code != http.StatusNoContent {
				t.Errorf("del wrong error code. expected %d got %d", http.StatusNoContent, w.Code)
			}

			f, ok := mustGetFavorite(t, srv, user, id)
			if ok || f == addedFavorite {
				t.Errorf("get retrieved deleted value")
			}
		})
	})

	t.Run("update", func(t *testing.T) {
		srv := newFavoritesService(newFavoritesStorage(), nil)
		addedFavorite := Favorite{AssetKey: "key", Description: "descr"}
		id := mustAddFavorite(t, srv, user, addedFavorite)

		t.Run("nonNumericID", func(t *testing.T) {
			w := callUpdate(srv, user, "hello", nil)
			if w.Code != http.StatusNotFound {
				t.Errorf("update accepted non numeric id")
			}
		})

		t.Run("malformedInput", func(t *testing.T) {
			w := callUpdate(srv, user, id, bytes.NewBufferString("Hello"))
			if w.Code != http.StatusBadRequest {
				t.Errorf("update accepted malformed input")
			}
		})

		t.Run("nonExistingFavorite", func(t *testing.T) {
			updateFavorite := Favorite{Description: "descrUpdated"}
			w := callUpdate(srv, user, invalidID, encodeJson(t, updateFavorite))
			if w.Code != http.StatusNotFound {
				t.Errorf("update accepted non existing favorite")
			}
		})

		t.Run("existingFavorite", func(t *testing.T) {
			// PUT input may not include the key
			updatedFavorite := Favorite{Description: "descrUpdated"}
			w := callUpdate(srv, user, id, encodeJson(t, updatedFavorite))
			if w.Code != http.StatusOK {
				t.Errorf("update error code. expected %d got %d", http.StatusOK, w.Code)
				t.FailNow()
			}

			expectedFavorite := addedFavorite
			expectedFavorite.ID, _ = strconv.Atoi(id)
			expectedFavorite.Description = updatedFavorite.Description // only this should change

			var returnedFavorite Favorite
			if err := json.NewDecoder(w.Body).Decode(&returnedFavorite); err != nil {
				t.Errorf("failed to parse update response: %v", err)
				t.FailNow()
			}
			if returnedFavorite != expectedFavorite {
				t.Errorf("update response error. expected %v got %v", expectedFavorite, returnedFavorite)
				t.FailNow()
			}

			favFromGet, ok := mustGetFavorite(t, srv, user, id)
			if !ok || favFromGet != expectedFavorite {
				t.Errorf("update did not made the updates, value after get not updated")
			}
		})
	})

	t.Run("getAll", func(t *testing.T) {
		silenceLogger(t)
		user := User{Email: "user@gwifav.dev", Password: "passwd"}

		// load test assets and start an asset server
		// The server has 20 assets of every kind
		var assets Assets
		fin, err := os.Open("testdata/test-assets.json")
		if err == nil {
			defer fin.Close()
			err = json.NewDecoder(fin).Decode(&assets)
		}
		if err != nil {
			t.Errorf("failed to load test assets: %v", err)
			t.FailNow()
		}
		as := &AssetsServer{assets, 0}
		mux := http.NewServeMux()
		mux.HandleFunc("GET /assets/insights", http.NotFound)
		mux.HandleFunc("GET /assets/insights/{$}", http.NotFound)
		mux.HandleFunc("GET /assets/insights/{key...}", as.insightsHandler)
		mux.HandleFunc("GET /assets/charts", http.NotFound)
		mux.HandleFunc("GET /assets/charts/{$}", http.NotFound)
		mux.HandleFunc("GET /assets/charts/{key...}", as.chartsHandler)
		mux.HandleFunc("GET /assets/audiences", http.NotFound)
		mux.HandleFunc("GET /assets/audiences/{$}", http.NotFound)
		mux.HandleFunc("GET /assets/audiences/{key...}", as.audiencesHandler)
		srv := httptest.NewServer(mux)
		srvURL, err := url.Parse(srv.URL)
		srvURL = srvURL.JoinPath("assets")
		if err != nil {
			t.Errorf("failed to parse server url %s: %v", srv.URL, err)
			t.FailNow()
		}
		t.Cleanup(func() { srv.Close() })

		// create a favorites handler with some favorites
		nAddedFavorites := 9
		fstore := newFavoritesStorage()
		var key string
		for i := 1; i <= nAddedFavorites/3; i++ {
			key = fmt.Sprintf("insights/%d", i)
			fstore.add(user.Email, Favorite{AssetKey: key, Description: "description of " + key})
			key = fmt.Sprintf("charts/%d", i)
			fstore.add(user.Email, Favorite{AssetKey: key, Description: "description of " + key})
			key = fmt.Sprintf("audiences/%d", i)
			fstore.add(user.Email, Favorite{AssetKey: key, Description: "description of " + key})
		}
		// add some more without assets, to check the server returns nulls and don't fail
		key = "insights/1000"
		fstore.add(user.Email, Favorite{AssetKey: key, Description: "description of " + key})
		key = "charts/1000"
		fstore.add(user.Email, Favorite{AssetKey: key, Description: "description of " + key})
		key = "audiences/1000"
		fstore.add(user.Email, Favorite{AssetKey: key, Description: "description of " + key})
		nAddedFavorites += 3

		af := &AssetsFetcher{
			cli: http.DefaultClient,
			srv: srvURL,
		}
		fs := &FavoritesService{fstore, af}
		ctx := context.WithValue(context.Background(), keyUser, user)
		ctx = context.WithValue(ctx, keyToken, "dummy")
		req := httptest.NewRequestWithContext(ctx, "GET", "/favorites/getAll", nil)
		w := httptest.NewRecorder()
		fs.getAll(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("response code %d", w.Code)
			t.FailNow()
		}

		var response GetAllFavoritesResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Errorf("failed to decode response: %v", err)
			t.FailNow()
		}

		if response.Total != nAddedFavorites {
			t.Errorf("response total mismatch. expected %d got %d", nAddedFavorites, response.Total)
		}
		if len(response.Favorites) != nAddedFavorites {
			t.Errorf("response data mismatch. expected %d got %d", nAddedFavorites, len(response.Favorites))
		}

		audience1 := response.Favorites[slices.IndexFunc(response.Favorites, func(f FavoriteWithAsset) bool {
			return f.AssetKey == "audiences/1"
		})]
		if !audience1.IncludesAsset {
			t.Error("audience1 does not include asset but it should")
		}
		audience1000 := response.Favorites[slices.IndexFunc(response.Favorites, func(f FavoriteWithAsset) bool {
			return f.AssetKey == "audiences/1000"
		})]
		if audience1000.IncludesAsset {
			t.Error("audience1000 includes asset but it shouldn't")
		}
	})
}

//
// a simple http client for the service
//

func encodeJson(t *testing.T, fav Favorite) *bytes.Buffer {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(fav); err != nil {
		t.Errorf("failed to encode Favorite: %v", err)
		t.FailNow()
	}
	return body
}

func callAdd(srv *FavoritesService, user User, body io.Reader) *httptest.ResponseRecorder {
	ctx := context.WithValue(context.Background(), keyUser, user)
	req := httptest.NewRequestWithContext(ctx, "POST", "/favorites", body)
	w := httptest.NewRecorder()
	srv.add(w, req)
	return w
}

func callGet(srv *FavoritesService, user User, id string) *httptest.ResponseRecorder {
	ctx := context.WithValue(context.Background(), keyUser, user)
	u := fmt.Sprintf("/favorites/%v", id)
	req := httptest.NewRequestWithContext(ctx, "GET", u, nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	srv.get(w, req)
	return w
}

func callDelete(srv *FavoritesService, user User, id string) *httptest.ResponseRecorder {
	ctx := context.WithValue(context.Background(), keyUser, user)
	u := fmt.Sprintf("/favorites/%v", id)
	req := httptest.NewRequestWithContext(ctx, "DELETE", u, nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	srv.del(w, req)
	return w
}

func callUpdate(srv *FavoritesService, user User, id string, body io.Reader) *httptest.ResponseRecorder {
	ctx := context.WithValue(context.Background(), keyUser, user)
	u := fmt.Sprintf("/favorites/%v", id)
	req := httptest.NewRequestWithContext(ctx, "PUT", u, body)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	srv.update(w, req)
	return w
}

//
// utilites for testing
//

func mustAddFavorite(t *testing.T, srv *FavoritesService, user User, fav Favorite) string {
	w := callAdd(srv, user, encodeJson(t, fav))
	if w.Code != http.StatusCreated {
		t.Errorf("add failed")
		t.FailNow()
	}
	loc := w.Header().Get("Location")
	toks := strings.Split(loc, "/")
	return toks[2]
}

func mustGetFavorite(t *testing.T, srv *FavoritesService, user User, id string) (Favorite, bool) {
	w := callGet(srv, user, id)
	if w.Code == http.StatusNotFound {
		return Favorite{}, false
	}
	if w.Code == http.StatusOK {
		var fav Favorite
		if err := json.NewDecoder(w.Body).Decode(&fav); err != nil {
			t.Errorf("get failed: %v", err)
			t.FailNow()
		}
		return fav, true
	}
	t.Errorf("get returned unexpected code: %d", w.Code)
	t.FailNow()
	return Favorite{}, false // to silence the compiler
}

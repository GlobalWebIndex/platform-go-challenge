package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type response struct {
	Message string `json:"message"`
}

func decodeJSON(data []byte, v any) error {
	var unmarshalErr *json.UnmarshalTypeError
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		if errors.As(err, &unmarshalErr) {
			return fmt.Errorf("Invalid JSON format: %w", err)
		} else {
			return fmt.Errorf("Invalid JSON format: %w", err)
		}
	}
	return nil
}

func TestGetFavourites4User(t *testing.T) {
	client := &http.Client{}

	t. Run("happy path", func(t *testing.T) {
		url := "http://0.0.0.0:8082/users/user42/favourites?page=1&pageSize=10"
		resp, err := client.Get(url)
		require.NoError(t, err, "failed to get user42 favourites")
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("user not in favourites collection", func(t *testing.T) {
		url := "http://0.0.0.0:8082/users/user43/favourites?page=1&pageSize=10"
		resp, err := client.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("wrong pagination param format", func(t *testing.T) {
		url := "http://0.0.0.0:8082/users/user43/favourites?page=hello&pageSize=10"
		resp, err := client.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("page+pagiSize exceeds favourites count", func(t *testing.T) {
		url := "http://0.0.0.0:8082/users/user42/favourites?page=2&pageSize=4"
		resp, err := client.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

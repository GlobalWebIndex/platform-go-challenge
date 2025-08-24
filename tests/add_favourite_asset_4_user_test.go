package main

import (
	"bytes"
	"encoding/json"
	"gwiEng/internal/api/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//word = {"msg":"Asset with id = d9a972b5-2783-4b72-9cec-5f3031d61ab5 registered successfully as favourite "}

func TestAddFavourite(t *testing.T) {
	client := &http.Client{}

	t.Run("happy path", func(t *testing.T){
		url := "http://0.0.0.0:8082/users/user42/favourites/d9a972b5-2783-4b72-9cec-5f3031d61ab5"
		contentType := "application/json; charset=UTF-8"
		addFavBody := models.AddFavouriteBody{
			AssetType: "insight",
			Description: "A bried description",
		}
		addFavJSON, err := json.Marshal(addFavBody)		
		require.NoError(t, err)
		resp, err := client.Post(url, contentType, bytes.NewBuffer(addFavJSON))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("favourite assetID+userID exists", func(t *testing.T){
		url := "http://0.0.0.0:8082/users/user42/favourites/d9a972b5-2783-4b72-9cec-5f3031d61ab5"
		contentType := "application/json; charset=UTF-8"
		addFavBody := models.AddFavouriteBody{
			AssetType: "insight",
			Description: "A bried description",
		}
		addFavJSON, err := json.Marshal(addFavBody)		
		require.NoError(t, err)
		resp, err := client.Post(url, contentType, bytes.NewBuffer(addFavJSON))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("userID does not exist", func(t *testing.T){
		url := "http://0.0.0.0:8082/users/user43/favourites/d9a972b5-2783-4b72-9cec-5f3031d61ab5"
		contentType := "application/json; charset=UTF-8"
		addFavBody := models.AddFavouriteBody{
			AssetType: "insight",
			Description: "A bried description",
		}
		addFavJSON, err := json.Marshal(addFavBody)		
		require.NoError(t, err)
		resp, err := client.Post(url, contentType, bytes.NewBuffer(addFavJSON))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("assetID does not exist", func(t *testing.T){
		url := "http://0.0.0.0:8082/users/user42/favourites/d9a972b5-2783-4b72-9cec-5f3031d61ab6"
		contentType := "application/json; charset=UTF-8"
		addFavBody := models.AddFavouriteBody{
			AssetType: "insight",
			Description: "A bried description",
		}
		addFavJSON, err := json.Marshal(addFavBody)		
		require.NoError(t, err)
		resp, err := client.Post(url, contentType, bytes.NewBuffer(addFavJSON))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	type wrongReqBody struct {
		Asset string
	}

	t.Run("invalid request body", func(t *testing.T){
		url := "http://0.0.0.0:8082/users/user42/favourites/d9a972b5-2783-4b72-9cec-5f3031d61ab5"
		contentType := "application/json; charset=UTF-8"
		addFavBody := wrongReqBody{
			Asset: "insight",
		}
		addFavJSON, err := json.Marshal(addFavBody)		
		require.NoError(t, err)
		resp, err := client.Post(url, contentType, bytes.NewBuffer(addFavJSON))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})




} 












package gwi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ResponseAPI represents the data structure needed to create a response
type ResponseAPI struct {
	Status      int                     `json:"status"`
	Description string                  `json:"description,omitempty"`
	Success     bool                    `json:"success"`
	Length      int                     `json:"length,omitempty"`
	Data        map[AssetType][]Asseter `json:"data,omitempty"`
}

// response sets the params to generate a JSON response
func response(w http.ResponseWriter, ra ResponseAPI) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ra.Status)

	log.Printf("Response => %d", ra.Status)

	json.NewEncoder(w).Encode(ra)
}

// responseError returns a 5XX response
func responseError(w http.ResponseWriter, description string) {
	log.Println(description)

	ra := ResponseAPI{
		Status:      http.StatusInternalServerError,
		Description: description,
		Success:     false,
	}
	response(w, ra)
}

// responseNoContent returns a no content response
func responseNoContent(w http.ResponseWriter, description string) {
	log.Println(description)

	ra := ResponseAPI{
		Status:      http.StatusNoContent,
		Description: description,
		Success:     false,
	}
	response(w, ra)
}

// responseUnprocessable calls response function with proper data to generate a Unprocessable Entity response
func responseUnprocessable(w http.ResponseWriter, message string) {
	ra := ResponseAPI{
		Status:      http.StatusUnprocessableEntity,
		Description: fmt.Sprintf("%s input data not valid", message),
		Success:     false,
	}

	response(w, ra)
}

// responseOk calls response function with proper data to generate an OK response
func responseOk(w http.ResponseWriter, success bool) {
	ra := ResponseAPI{
		Status:  http.StatusOK,
		Success: success,
	}
	response(w, ra)
}

// responseList returns a List of favourites
func responseList(w http.ResponseWriter, success bool, assets map[AssetType][]Asseter) {
	ra := ResponseAPI{
		Status:  http.StatusOK,
		Success: success,
		Length:  len(assets),
		Data:    assets,
	}
	response(w, ra)
}

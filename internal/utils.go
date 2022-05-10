package gwi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetUuid retrieves the uuid from the url path using mux.Vars function
func GetUuid(req *http.Request) string {
	log.Println("Extracting uuid value from path")
	vars := mux.Vars(req)
	return vars["userid"]
}

// ValidateUuid validates the uuid format
func ValidateUuid(input string) bool {
	log.Println("Validating uuid value")

	_, err := uuid.Parse(input)
	return err == nil
}

// Decode decodes body request to an structure represented by the input type
func Decode(r *http.Request, input interface{}) (interface{}, error) {
	log.Println("Decoding input")

	switch input.(type) {
	case Asset:
		output := Asset{}
		if err := json.NewDecoder(r.Body).Decode(&output); err != nil {
			log.Printf("Decode error => Asset => %v", err)
			return nil, err
		}
		defer r.Body.Close()
		return output, nil
	}
	return nil, nil
}

// PrettyPrint is a helper function to print structs
func PrettyPrint(input interface{}) {
	empJSON, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("PrettyPrint output \n %s\n", string(empJSON))
}

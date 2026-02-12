// Command swaggergen generates OpenAPI 3.0 specification files (JSON and YAML)
// for the Platform Go Challenge Favourites API and writes them to the api/ directory.
//
// Usage:
//
//	go run ./tools/swaggergen
//
// # For Contributors
//
// When you modify the API (add/change endpoints, request/response schemas, etc.),
// update this file to keep the swagger spec in sync:
//
//  1. Endpoints: Edit buildPaths() to add/modify path items and operations
//  2. Schemas: Edit buildSchemas() to add/modify request/response types
//  3. Regenerate: Run `go run ./tools/swaggergen` from the project root
//  4. Verify: Check api/swagger.yaml and api/swagger.json for correctness
//
// Helper functions:
//   - errContent(): Returns standard error response content (reuse for error responses)
//   - assetIDParam(): Returns the {assetID} path parameter definition
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// ---------------------------------------------------------------------------
// Lightweight OpenAPI 3.0 types
// ---------------------------------------------------------------------------

type OpenAPI struct {
	OpenAPI    string               `json:"openapi"              yaml:"openapi"`
	Info       Info                 `json:"info"                 yaml:"info"`
	Paths      map[string]*PathItem `json:"paths"                yaml:"paths"`
	Components Components           `json:"components"           yaml:"components"`
}

type Info struct {
	Title       string `json:"title"       yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Version     string `json:"version"     yaml:"version"`
}

type PathItem struct {
	Get    *Operation `json:"get,omitempty"    yaml:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"   yaml:"post,omitempty"`
	Patch  *Operation `json:"patch,omitempty"  yaml:"patch,omitempty"`
	Delete *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
}

type Operation struct {
	Tags        []string              `json:"tags"                  yaml:"tags"`
	Summary     string                `json:"summary"               yaml:"summary"`
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	OperationID string                `json:"operationId"           yaml:"operationId"`
	Security    []map[string][]string `json:"security,omitempty"    yaml:"security,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"  yaml:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"             yaml:"responses"`
}

type Parameter struct {
	Name        string `json:"name"        yaml:"name"`
	In          string `json:"in"          yaml:"in"`
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"required"    yaml:"required"`
	Schema      Schema `json:"schema"      yaml:"schema"`
}

type RequestBody struct {
	Required    bool                 `json:"required"              yaml:"required"`
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]MediaType `json:"content"               yaml:"content"`
}

type MediaType struct {
	Schema Schema `json:"schema" yaml:"schema"`
}

type Response struct {
	Description string               `json:"description"       yaml:"description"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

type Schema struct {
	Type                 string            `json:"type,omitempty"                 yaml:"type,omitempty"`
	Format               string            `json:"format,omitempty"               yaml:"format,omitempty"`
	Description          string            `json:"description,omitempty"          yaml:"description,omitempty"`
	Properties           map[string]Schema `json:"properties,omitempty"           yaml:"properties,omitempty"`
	Items                *Schema           `json:"items,omitempty"                yaml:"items,omitempty"`
	Required             []string          `json:"required,omitempty"             yaml:"required,omitempty"`
	Enum                 []string          `json:"enum,omitempty"                 yaml:"enum,omitempty"`
	Ref                  string            `json:"$ref,omitempty"                 yaml:"$ref,omitempty"`
	AdditionalProperties *Schema           `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	OneOf                []Schema          `json:"oneOf,omitempty"                yaml:"oneOf,omitempty"`
	Example              any               `json:"example,omitempty"              yaml:"example,omitempty"`
}

type Components struct {
	Schemas         map[string]Schema         `json:"schemas"         yaml:"schemas"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes" yaml:"securitySchemes"`
}

type SecurityScheme struct {
	Type         string `json:"type"         yaml:"type"`
	Scheme       string `json:"scheme"       yaml:"scheme"`
	BearerFormat string `json:"bearerFormat" yaml:"bearerFormat"`
	Description  string `json:"description"  yaml:"description"`
}

// ---------------------------------------------------------------------------
// Spec builder
// ---------------------------------------------------------------------------

func buildSpec() OpenAPI {
	bearerAuth := []map[string][]string{{"BearerAuth": {}}}

	return OpenAPI{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:       "Platform Go Challenge - Favourites API",
			Description: "REST API for managing user favourite assets (charts, insights, audiences).",
			Version:     "1.0.0",
		},
		Paths: buildPaths(bearerAuth),
		Components: Components{
			Schemas:         buildSchemas(),
			SecuritySchemes: buildSecuritySchemes(),
		},
	}
}

func buildPaths(bearerAuth []map[string][]string) map[string]*PathItem {
	return map[string]*PathItem{
		"/api/v1/favourites": {
			Get: &Operation{
				Tags:        []string{"Favourites"},
				Summary:     "List user favourites",
				Description: "Returns all favourite assets for the authenticated user.",
				OperationID: "getUserFavourites",
				Security:    bearerAuth,
				Responses: map[string]Response{
					"200": {
						Description: "A list of favourite assets",
						Content: map[string]MediaType{
							"application/json": {Schema: Schema{
								Type:  "array",
								Items: &Schema{Ref: "#/components/schemas/FavouriteAsset"},
							}},
						},
					},
					"401": {Description: "Unauthorized - missing or invalid JWT"},
					"406": {Description: "Not Acceptable - Accept header must include application/json", Content: errContent()},
					"500": {Description: "Internal server error", Content: errContent()},
				},
			},
			Post: &Operation{
				Tags:        []string{"Favourites"},
				Summary:     "Add a favourite",
				Description: "Adds a new asset to the authenticated user's favourites.",
				OperationID: "addUserFavourite",
				Security:    bearerAuth,
				RequestBody: &RequestBody{
					Required:    true,
					Description: "Asset to favourite",
					Content: map[string]MediaType{
						"application/json": {Schema: Schema{Ref: "#/components/schemas/AddFavouriteRequest"}},
					},
				},
				Responses: map[string]Response{
					"201": {
						Description: "Favourite added",
						Content: map[string]MediaType{
							"application/json": {Schema: Schema{Ref: "#/components/schemas/SuccessMessage"}},
						},
					},
					"400": {Description: "Invalid request body or validation error", Content: errContent()},
					"401": {Description: "Unauthorized"},
					"406": {Description: "Not Acceptable - Accept header must include application/json", Content: errContent()},
					"409": {Description: "Favourite already exists", Content: errContent()},
					"415": {Description: "Unsupported Media Type - Content-Type must be application/json", Content: errContent()},
					"500": {Description: "Internal server error", Content: errContent()},
				},
			},
		},
		"/api/v1/favourites/{assetID}": {
			Patch: &Operation{
				Tags:        []string{"Favourites"},
				Summary:     "Update favourite description",
				Description: "Updates the description of an existing favourite asset.",
				OperationID: "updateUserFavourite",
				Security:    bearerAuth,
                Parameters:  []Parameter{assetIDParam()},
				RequestBody: &RequestBody{
					Required: true,
					Content: map[string]MediaType{
						"application/json": {Schema: Schema{Ref: "#/components/schemas/UpdateDescriptionRequest"}},
					},
				},
				Responses: map[string]Response{
					"200": {
						Description: "Description updated",
						Content: map[string]MediaType{
							"application/json": {Schema: Schema{Ref: "#/components/schemas/SuccessMessage"}},
						},
					},
					"400": {Description: "Invalid request body or validation error", Content: errContent()},
					"401": {Description: "Unauthorized"},
					"404": {Description: "Favourite not found", Content: errContent()},
					"406": {Description: "Not Acceptable - Accept header must include application/json", Content: errContent()},
					"415": {Description: "Unsupported Media Type - Content-Type must be application/json", Content: errContent()},
					"500": {Description: "Internal server error", Content: errContent()},
				},
			},
			Delete: &Operation{
				Tags:        []string{"Favourites"},
				Summary:     "Remove a favourite",
				Description: "Removes an asset from the authenticated user's favourites.",
				OperationID: "removeUserFavourite",
				Security:    bearerAuth,
                Parameters:  []Parameter{assetIDParam()},
				Responses: map[string]Response{
					"200": {
						Description: "Favourite removed",
						Content: map[string]MediaType{
							"application/json": {Schema: Schema{Ref: "#/components/schemas/SuccessMessage"}},
						},
					},
					"400": {Description: "Missing asset ID", Content: errContent()},
					"401": {Description: "Unauthorized"},
					"404": {Description: "Favourite not found", Content: errContent()},
					"406": {Description: "Not Acceptable - Accept header must include application/json", Content: errContent()},
					"500": {Description: "Internal server error", Content: errContent()},
				},
			},
		},
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func assetIDParam() Parameter {
	return Parameter{
		Name:        "assetID",
		In:          "path",
		Description: "Unique identifier of the favourite asset",
		Required:    true,
		Schema:      Schema{Type: "string"},
	}
}



func errContent() map[string]MediaType {
	return map[string]MediaType{
		"application/json": {Schema: Schema{Ref: "#/components/schemas/ErrorResponse"}},
	}
}

func buildSecuritySchemes() map[string]SecurityScheme {
	return map[string]SecurityScheme{
		"BearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "JWT token with a 'sub' claim identifying the user.",
		},
	}
}

func buildSchemas() map[string]Schema {
	return map[string]Schema{
		"ErrorResponse": {
			Type: "object",
			Properties: map[string]Schema{
				"error": {Type: "string", Description: "Human-readable error message"},
			},
			Required: []string{"error"},
		},
		"SuccessMessage": {
			Type: "object",
			Properties: map[string]Schema{
				"message": {Type: "string", Description: "Success message"},
			},
			Required: []string{"message"},
		},
		"AddFavouriteRequest": {
			Type:        "object",
			Description: "Payload for adding a favourite asset. The asset_data shape depends on asset_type.",
			Properties: map[string]Schema{
				"asset_type": {
					Type: "string",
					Enum: []string{"chart", "insight", "audience"},
					Description: "Type of asset being favourited",
				},
				"description": {
					Type:        "string",
					Description: "Optional description for the favourite (max 255 chars)",
				},
				"asset_data": {
					Description: "Asset payload - one of Chart, Insight or Audience",
					OneOf: []Schema{
						{Ref: "#/components/schemas/Chart"},
						{Ref: "#/components/schemas/Insight"},
						{Ref: "#/components/schemas/Audience"},
					},
				},
			},
			Required: []string{"asset_type", "asset_data"},
		},
		"UpdateDescriptionRequest": {
			Type: "object",
			Properties: map[string]Schema{
				"description": {Type: "string", Description: "New description (max 255 chars)"},
			},
			Required: []string{"description"},
		},
		"FavouriteAsset": {
			Type:        "object",
			Description: "A user's favourited asset with metadata.",
			Properties: map[string]Schema{
				"id":          {Type: "string"},
				"user_id":     {Type: "string"},
				"asset_type":  {Type: "string", Enum: []string{"chart", "insight", "audience"}},
				"description": {Type: "string"},
				"created_at":  {Type: "string", Format: "date-time"},
				"updated_at":  {Type: "string", Format: "date-time"},
				"data": {
					Description: "The full asset object",
					OneOf: []Schema{
						{Ref: "#/components/schemas/Chart"},
						{Ref: "#/components/schemas/Insight"},
						{Ref: "#/components/schemas/Audience"},
					},
				},
			},
			Required: []string{"id", "user_id", "asset_type", "created_at", "updated_at", "data"},
		},
		"Chart": {
			Type:        "object",
			Description: "A chart asset.",
			Properties: map[string]Schema{
				"id":           {Type: "string"},
				"title":        {Type: "string"},
				"x_axis_title": {Type: "string"},
				"y_axis_title": {Type: "string"},
				"data": {
					Type:                 "object",
					AdditionalProperties: &Schema{},
					Description:          "Arbitrary chart data points",
				},
			},
			Required: []string{"id", "title", "x_axis_title", "y_axis_title"},
		},
		"Insight": {
			Type:        "object",
			Description: "An insight asset.",
			Properties: map[string]Schema{
				"id":   {Type: "string"},
				"text": {Type: "string"},
			},
			Required: []string{"id", "text"},
		},
		"Audience": {
			Type:        "object",
			Description: "An audience segment asset.",
			Properties: map[string]Schema{
				"id": {Type: "string"},
				"gender": {
					Type:  "array",
					Items: &Schema{Type: "string", Enum: []string{"Male", "Female"}},
				},
				"birth_country": {
					Type:  "array",
					Items: &Schema{Type: "string"},
				},
				"age_groups": {
					Type:  "array",
					Items: &Schema{Type: "string", Enum: []string{"18-24", "25-34", "35-44", "45-54", "55+"}},
				},
				"social_media_hours_daily": {
					Type: "string",
					Enum: []string{"0-1", "1-3", "3-5", "5+"},
				},
				"purchases_last_month": {
					Type:        "integer",
					Description: "Must be non-negative",
				},
			},
			Required: []string{"id"},
		},
	}
}

// ---------------------------------------------------------------------------
// File writers
// ---------------------------------------------------------------------------

func writeJSON(spec OpenAPI, path string) error {
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0644)
}

func writeYAML(spec OpenAPI, path string) error {
	data, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("marshal YAML: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func main() {
	_, src, _, _ := runtime.Caller(0)
	outDir := filepath.Join(filepath.Join(filepath.Dir(src), "..", ".."), "api")

	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create api/ directory: %v\n", err)
		os.Exit(1)
	}

	spec := buildSpec()

	jsonPath := filepath.Join(outDir, "swagger.json")
	if err := writeJSON(spec, jsonPath); err != nil {
		fmt.Fprintf(os.Stderr, "error writing JSON: %v\n", err)
		os.Exit(1)
	}

	yamlPath := filepath.Join(outDir, "swagger.yaml")
	if err := writeYAML(spec, yamlPath); err != nil {
		fmt.Fprintf(os.Stderr, "error writing YAML: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Swagger specs generated:\n  %s\n  %s\n", jsonPath, yamlPath)
}

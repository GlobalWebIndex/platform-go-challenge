package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

const docsHTML = `<!doctype html>
<html>
<head>
  <title>API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: "/openapi/swagger.json",
      dom_id: "#swagger-ui"
    });
  </script>
</body>
</html>`

func serveDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(docsHTML))
}

func serveSwaggerYAML(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("swagger.yaml")
	if err != nil {
		http.Error(w, fmt.Sprintf("read swagger: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func serveSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("swagger.yaml")
	if err != nil {
		http.Error(w, fmt.Sprintf("read swagger: %v", err), http.StatusInternalServerError)
		return
	}
	var v any
	if err := yaml.Unmarshal(data, &v); err != nil {
		http.Error(w, fmt.Sprintf("parse swagger: %v", err), http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		http.Error(w, fmt.Sprintf("encode swagger: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

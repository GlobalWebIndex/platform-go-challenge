package handlers

import (
    "fmt"
    "net/http"
)


// HomeHandler handles the home page request
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    http.ServeFile(w, r, "html/index.html")
}


// HomeHandler handles the home page request
func DocsHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/docs" {
        http.NotFound(w, r)
        return
    }
    http.ServeFile(w, r, "html/api_docs.html")
}


// PingHandler handles the /ping endpoint request
func PingHandler(w http.ResponseWriter, r *http.Request) {
    response := fmt.Sprintf("OK")
    fmt.Fprint(w, response)
}
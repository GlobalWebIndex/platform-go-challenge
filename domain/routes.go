package domain

import "net/http"

type RouteDefinition struct {
	Methods     []string
	HandlerFunc http.HandlerFunc
}

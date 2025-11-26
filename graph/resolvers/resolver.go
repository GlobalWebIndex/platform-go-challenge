package resolvers

import (
	"platform-go-challenge/graph"

	"gorm.io/gorm"
)

type Resolver struct {
	DB *gorm.DB
}

// Ensure Resolver implements the graph.ResolverRoot interface
var _ graph.ResolverRoot = &Resolver{}

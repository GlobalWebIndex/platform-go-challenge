package mongo_test

import (
	"context"
	mongorepo "platform-go-challenge/internal/infra/repository/mongo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	expErrContainsStr := "error validating uri: username required if URI contains user info"
	tests := map[string]struct {
		cfg            mongorepo.Config
		expErrContains *string
	}{
		"should return error on mongo client config error": {
			cfg:            mongorepo.Config{},
			expErrContains: &expErrContainsStr,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := mongorepo.New(context.TODO(), tt.cfg)
			if tt.expErrContains != nil {
				assert.ErrorContains(t, err, *tt.expErrContains)
			}
		})
	}
}

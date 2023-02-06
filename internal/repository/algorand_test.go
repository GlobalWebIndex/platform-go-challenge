package repository_test

import (
	"os"
	"ownify_api/internal/domain"
	"ownify_api/internal/repository"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestNewClient(t *testing.T) {
	client, index, err := repository.NewClient(domain.TestNet)
	require.Nil(t, err)
	clientHealth := client.HealthCheck()
	require.NotNil(t, clientHealth)
	indexerHealth := index.HealthCheck()
	require.NotNil(t, indexerHealth)
}

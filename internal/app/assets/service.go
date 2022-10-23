package assets

import "context"

type Service struct {
	repository Repository
}

func NewAssetService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context) (Assets, error) {
	return s.repository.List(ctx)
}

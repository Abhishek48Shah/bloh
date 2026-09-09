package blog

import "context"

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
}
func (s *Service) FindById(ctx context.Context, id int) (*Blog, error) {
	return nil, nil
}

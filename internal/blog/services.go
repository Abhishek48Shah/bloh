package blog

import (
	"context"
)

type Service struct {
	storage *Storage
}

const (
	defaultLimit = 10
	maxLimit     = 50
)

func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
}
func (s *Service) GetBlogById(ctx context.Context, id int) (*Blog, error) {
	return s.storage.FindBlogById(ctx, id)
}
func (s *Service) GetUserBlogs(ctx context.Context, userID int, pagination *Pagination) ([]Blog, error) {
	if pagination.limit < 0 {
		pagination.limit = defaultLimit
	}
	if pagination.limit > maxLimit {
		pagination.limit = 50
	}
	blogs, err := s.storage.FindUserBlogs(ctx, userID, pagination)
	if err != nil {
		return nil, err
	}
	pagination.curosr = blogs[len(blogs)-1].CreatedAt
	return blogs, nil
}
func (s *Service) CreateBlog(ctx context.Context, blog *CreateBlogInput, userID int) error {
	return s.storage.CreateBlog(ctx, blog, userID)
}
func (s *Service) DeleteBlog(ctx context.Context, blogID, userID int) error {
	return s.storage.Delete(ctx, blogID, userID)
}
func (s *Service) UpdateBlog(ctx context.Context, input *UpdateBlogInput, blogID, userID int) error {
	return s.storage.Update(ctx, input, blogID, userID)
}

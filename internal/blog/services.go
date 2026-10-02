package blog

import (
	"context"
)

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
}
func (s *Service) GetBlogById(ctx context.Context, id int) (*Blog, error) {
	return s.storage.FindBlogById(ctx, id)
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

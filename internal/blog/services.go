package blog

import (
	"context"
	"fmt"
)

type BlogServices struct {
	storage *BlogStorage
}

func NewBlogServices(storage *BlogStorage) *BlogServices {
	return &BlogServices{storage: storage}
}
func (s *BlogServices) GetBlogById(ctx context.Context, blogID int) (*Blog, error) {
	var blog Blog
	err := s.storage.FindBlogById(ctx, blogID, &blog)
	if err != nil {
		return nil, fmt.Errorf("services/GetBlogById: %w", err)
	}
	return &blog, nil
}
func (s *BlogServices) CreateBlog(ctx context.Context, userID int, data *UserInput) error {
	parsedData := Parser(data)
	return nil
}

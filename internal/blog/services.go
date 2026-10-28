package blog

import (
	"context"

	"github.com/Abhishek48Shah/blog/internal/pagination"
)

type Service struct {
	storage *Storage
}

const limit = 10

func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
}
func (s *Service) GetBlogById(ctx context.Context, id int) (*Blog, error) {
	return s.storage.FindBlogById(ctx, id)
}
func (s *Service) GetUserBlogs(ctx context.Context, userID int, cursor *pagination.Cursor) (*BlogPage, error) {
	blogs, err := s.storage.FindUserBlogs(ctx, userID, cursor, limit)
	if err != nil {
		return nil, err
	}
	resp := &BlogPage{Blogs: blogs}
	if len(blogs) > limit {
		blogs = blogs[:limit]
		last := blogs[len(blogs)-1]
		resp.Blogs = blogs
		resp.NextCursor = &pagination.Cursor{CreatedAt: &last.CreatedAt, BlogID: &last.Id}
	}
	return resp, nil
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

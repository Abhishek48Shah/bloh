package blog

import (
	"context"
	"errors"
	"time"

	"github.com/Abhishek48Shah/blog/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	pool  *pgxpool.Pool
	cache *redis.Client
}

func NewStorage(pool *pgxpool.Pool, cache *redis.Client) *Storage {
	return &Storage{pool: pool, cache: cache}
}
func (s *Storage) FindBlogById(ctx context.Context, id int, blog *Blog) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `SELECT blog_id,title,content,created_at,updated_at FROM blogs WHERE blog_id=$1`
	err := s.pool.QueryRow(ctx, query, id).Scan(&blog.Id, &blog.Title, &blog.Content, &blog.CreateAt, &blog.UpdateAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return util.NotFounErr("blog not found", err)
		case errors.Is(err, context.DeadlineExceeded):
			return util.TimeoutErr("request time out", err)
		default:
			return util.InternalServerErr("internal server error", err)
		}
	}
	return nil
}
func (s *Storage) CreateBlog(ctx context.Context) {

}

package blog

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
)

type Storage struct {
	pool  *pgxpool.Pool
	cache *redis.Client
}

var (
	ErrNotFound = errors.New("blog not found")
	ErrConflict = errors.New("blog already exists")
)

func NewStorage(pool *pgxpool.Pool, cache *redis.Client) *Storage {
	return &Storage{pool: pool, cache: cache}
}
func (s *Storage) FindBlogById(ctx context.Context, blogID int) (*Blog, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `SELECT blog_id,title,content,created_at,updated_at FROM blogs WHERE blog_id=$1`
	var blog *Blog
	err := s.pool.QueryRow(ctx, query, blogID).Scan(&blog.Id, &blog.Title, &blog.Content, &blog.CreatedAt, &blog.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, fmt.Errorf("find blog %d: %w", blogID, ErrNotFound)
		case errors.Is(err, context.DeadlineExceeded):
			return nil, fmt.Errorf("find blog %d: %w", blogID, err)
		case errors.Is(err, context.Canceled):
			return nil, fmt.Errorf("find blog %d: %w", err)
		default:
			return nil, fmt.Errorf("find blog %d: %w", blogID, err)
		}
	}
	return blog, nil
}
func (s *Storage) CreateBlog(ctx context.Context, blog *CreateBlogInput, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `INSERT INTO blogs (user_id,title,content) VALUES ($1,$2,$3)`
	args := []any{userID, blog.Title, blog.Content}
	_, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("create blog: %w", err)
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("create blog: %w", err)
		case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
			return fmt.Errorf("create blog: %w", ErrConflict)
		default:
			return fmt.Errorf("create blog: %w", err)
		}
	}

	return nil
}
func (s *Storage) Delete(ctx context.Context, blogID, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `DELETE FROM blogs WHERE blog_id = $1 AND user_id = $2`
	result, err := s.pool.Exec(ctx, query, blogID, userID)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("delete blog: %w", err)
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("delete blog: %w", err)
		default:
			return fmt.Errorf("delete blog %d: %w", blogID, err)
		}

	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete blog %d: %w", blogID, ErrNotFound)
	}
	return nil
}

func (s *Storage) Update(ctx context.Context, blog *UpdateBlogInput, blogID, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `UPDATE blogs SET title=COALESCE($1,title),content=COALESCE($2,content) WHERE blog_id=$3 and user_id=$4`
	args := []any{blog.Title, blog.Content, blogID, userID}
	result, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("update blog %d: %w", blogID, err)
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("update blog %d: %w", blogID, err)
		default:
			return fmt.Errorf("update blog %d: %w", blogID, err)
		}
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("delete blog %d: %w", blogID, ErrNotFound)
	}
	return nil
}

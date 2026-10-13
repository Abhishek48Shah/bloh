package blog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"strconv"

	"github.com/Abhishek48Shah/blog/internal/middleware"
	"github.com/Abhishek48Shah/blog/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	service  *Service
	logger   *slog.Logger
	validate *validator.Validate
}

func NewHandler(service *Service, logger *slog.Logger, validate *validator.Validate) *Handler {
	return &Handler{service: service, logger: logger, validate: validate}
}

var maxBodySize int64 = 2097152

func (h *Handler) GetUserBlogsHandler(w http.ResponseWriter, r *http.Request) {
	payload := middleware.GetUserContext(r)
	blogs, err := h.service.GetUserBlogs(r.Context(), payload.ID)
	if err != nil {
		switch {
		default:
			util.WriteError(w, r, util.InternalServerErr("internal server error", err), h.logger)
		}
		return
	}
	// var header http.Header
	// header.Set("Cache-Controle")
	if err := util.WriteJson(w, blogs, http.StatusOK, nil); err != nil {
		util.WriteError(w, r, util.InternalServerErr("something went wrong", err), h.logger)
		return
	}
}
func (h *Handler) GetSingleBlogHandler(w http.ResponseWriter, r *http.Request) {
	blogID, err := util.ParseID(r)
	if err != nil {
		util.WriteError(w, r, util.BadRequestErr("invalid path parameter", fmt.Errorf("get blog: %w", err)), h.logger)
		return
	}
	blog, err := h.service.GetBlogById(r.Context(), blogID)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			util.WriteError(w, r, util.TimeoutErr("request time out", err), h.logger)
		case errors.Is(err, context.Canceled):
			return
		case errors.Is(err, ErrNotFound):
			util.WriteError(w, r, util.NotFounErr("blog not found", err), h.logger)
		case errors.Is(err, ErrConflict):
			util.WriteError(w, r, util.ConflictErr("blog already exists", err), h.logger)
		default:
			util.WriteError(w, r, util.InternalServerErr("internal server error", err), h.logger)
		}
		return
	}
	header := make(http.Header)
	header.Set("Cache-Control", "public")
	if err := util.WriteJson(w, blog, http.StatusOK, header); err != nil {
		util.WriteError(w, r, util.InternalServerErr("internal server error", err), h.logger)
		return
	}
}
func (h *Handler) CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		util.WriteError(w, r, util.UnsupportedMediaTypeErr("unsupported media type", err), h.logger)
		return
	}
	if mediaType != "application/json" {
		util.WriteError(w, r, util.UnsupportedMediaTypeErr("unsupported media type", err), h.logger)
		return
	}
	input := struct {
		Title   string `json:"title" validate:"required,min=21,max=200"`
		Content string `json:"content" validate:"required,min=50"`
	}{}
	if err := util.ReadJSON(w, r, maxBodySize, &input); err != nil {
		util.WriteError(w, r, util.BadRequestErr(err.Error(), err), h.logger)
		return
	}

	if err := validateBlog(h.validate, input); err != nil {
		util.WriteError(w, r, err, h.logger)
		return
	}
	blog := &CreateBlogInput{Title: input.Title, Content: input.Content}
	payload := middleware.GetUserContext(r)
	if err := h.service.CreateBlog(r.Context(), blog, payload.ID); err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			util.WriteError(w, r, util.TimeoutErr("request time out", err), h.logger)
		case errors.Is(err, context.Canceled):
			return
		case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
			util.WriteError(w, r, util.ConflictErr("blog already exist", err), h.logger)
		default:
			util.WriteError(w, r, util.InternalServerErr("internal server error", err), h.logger)
		}
		return
	}
	if err := util.WriteJson(w, "blog created successfully", http.StatusCreated, nil); err != nil {
		util.WriteError(w, r, util.InternalServerErr("internal server error", err), h.logger)
		return
	}
}
func (h *Handler) DeleteBlogHandler(w http.ResponseWriter, r *http.Request) {
	blogID, err := util.ParseID(r)
	if err != nil {
		util.WriteError(w, r, util.BadRequestErr("invalid path parameter", err), h.logger)
		return
	}
	payload := middleware.GetUserContext(r)
	if err := h.service.DeleteBlog(r.Context(), blogID, payload.ID); err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			util.WriteError(w, r, util.TimeoutErr("request time out", err), h.logger)
		case errors.Is(err, context.Canceled):
			return
		default:
			util.WriteError(w, r, util.InternalServerErr("internal server error", err), h.logger)
		}
		return

	}
	if err := util.WriteJson(w, "blog deleted successfully", http.StatusNoContent, nil); err != nil {
		fmt.Println(err)
		util.WriteError(w, r, util.InternalServerErr("internal server error", err), h.logger)
		return
	}
}
func (h *Handler) EditBlogHandler(w http.ResponseWriter, r *http.Request) {
	blogID := r.PathValue("id")
	parsedBlogId, err := strconv.Atoi(blogID)
	if err != nil {
		util.WriteError(w, r, util.BadRequestErr("invalid path parameter", err), h.logger)
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		util.WriteError(w, r, util.UnauthorizedErr("", fmt.Errorf("parse media type: %w", err)), h.logger)
		return
	}
	if mediaType != "application/json" {
		util.WriteError(w, r, util.UnauthorizedErr("", fmt.Errorf("parse media type: %w", err)), h.logger)
		return
	}
	input := struct {
		Title   *string `json:"title" validate:"omitnil,required,min=21,max=200"`
		Content *string `json:"content" validate:"omitnil,required,min=50,omitnil"`
	}{}

	if err := util.ReadJSON(w, r, maxBodySize, &input); err != nil {
		util.WriteError(w, r, util.BadRequestErr(err.Error(), fmt.Errorf("read json: %w", err)), h.logger)
	}
	if err := validateBlog(h.validate, input); err != nil {
		util.WriteError(w, r, err, h.logger)
		return
	}
	var blog UpdateBlogInput
	if input.Title != nil {
		blog.Title = input.Title
	}
	if input.Content != nil {
		blog.Content = input.Content
	}
	payload := middleware.GetUserContext(r)
	if err := h.service.UpdateBlog(r.Context(), &blog, parsedBlogId, payload.ID); err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			util.WriteError(w, r, util.TimeoutErr("request time out", err), h.logger)
		case errors.Is(err, context.Canceled):
			return
		default:
			util.WriteError(w, r, util.InternalServerErr("internal server error", err), h.logger)
		}
		return
	}
	if err := util.WriteJson(w, "blog update successfully", http.StatusOK, nil); err != nil {
		util.WriteError(w, r, util.InternalServerErr("internal server error", fmt.Errorf("write json: %w", err)), h.logger)
		return
	}
}

package blog

import (
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/Abhishek48Shah/blog/internal/middleware"
	"github.com/Abhishek48Shah/blog/pkg/util"
	"github.com/go-playground/validator/v10"
)

type BlogHandler struct {
	services  *BlogServices
	validator *validator.Validate
}

const maxBodySize = 1024

func NewBlogHandlers(services *BlogServices, validator *validator.Validate) *BlogHandler {
	return &BlogHandler{services: services, validator: validator}
}
func (h *BlogHandler) GetSingleBlog(w http.ResponseWriter, r *http.Request) {
	blogID := r.PathValue("id")
	parsedId, err := strconv.Atoi(blogID)
	if err != nil {
		err = fmt.Errorf("handler/GetSingleBlog: id parse: %w", err)
		util.BadRequest(w, r, err)
		return
	}
	if parsedId < 0 {
		util.BadRequest(w, r, fmt.Errorf("handler/GetSingleBlog: invalid parameter format"))
		return
	}
	ctx := r.Context()
	blog, err := h.services.GetBlogById(ctx, parsedId)
	if err != nil {
		util.InternalServerError(w, r, fmt.Errorf("handler/GetSingleBlog: %w", err))
		return
	}
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	maxAge := 7 * 24 * time.Hour
	headers.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge.Seconds()))
	err = util.WriteJSON(w, blog, http.StatusOK, headers)
	if err != nil {
		util.InternalServerError(w, r, fmt.Errorf("handle/GetSingleBlog: writeJSON: %w", err))
		return
	}
}
func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		util.UnsupportedMediaType(w, r)
		return
	}
	ctx := r.Context()
	userID, ok := middleware.GetUserId(ctx)
	if !ok {
		util.UnauthorizedError(w, r, fmt.Errorf("internal/blog/handler: user ID not found in context"), "unauthroized access")
		return
	}
	var userInput *UserInput
	err = util.ReadJSON(w, r, maxBodySize, &userInput)
	if err != nil {
		util.BadRequest(w, r, err)
		return
	}
	err = h.validator.Struct(userInput)
	if err != nil {
		return
	}
	err = h.services.CreateBlog(ctx, userID, userInput)
	if err != nil {
		return
	}
	err = util.WriteJSON(w, "blog created successfully", http.StatusCreated, nil)
	if err != nil {
		util.InternalServerError(w, r, err)
	}
}

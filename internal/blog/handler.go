package blog

import (
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/Abhishek48Shah/blog/pkg/util"
)

type BlogHandler struct {
	services *BlogServices
}

func NewBlogHandlers(services *BlogServices) *BlogHandler {
	return &BlogHandler{services: services}
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
}

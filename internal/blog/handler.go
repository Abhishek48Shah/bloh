package blog

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Abhishek48Shah/blog/internal/util"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service  *Service
	logger   *slog.Logger
	validate *validator.Validate
}

func NewHandler(service *Service, logger *slog.Logger, validate *validator.Validate) *Handler {
	return &Handler{service: service, logger: logger, validate: validate}
}
func (h *Handler) GetSingleBlog(w http.ResponseWriter, r *http.Request) {
	blogID := r.PathValue("id")
	err := h.validate.Var(blogID, "number")
	if err != nil {
		util.WriteError(w, r, util.BadRequestErr("invalid path parameter", err), h.logger)
		return
	}
	parseBlogId, err := strconv.Atoi(blogID)
	if err != nil {
		util.WriteError(w, r, util.BadRequestErr("invalid path parameter", err), h.logger)
		return
	}
	blog, err := h.service.FindById(r.Context(), parseBlogId)
	if err != nil {
		util.WriteError(w, r, err, h.logger)
	}
	var header http.Header
	header.Set("Content-Type", "application/json")
	header.Set("Cache-Control", "public")
	if err := util.WriteJson(w, blog, http.StatusAccepted, header); err != nil {
		util.WriteError(w, r, util.InternalServerErr("internal server error", err), h.logger)
		return
	}
}

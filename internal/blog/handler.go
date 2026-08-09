package blog

import (
	"net/http"
	"strconv"
)

type BlogHandler struct {
	services *BlogServices
}

func LoadHandlers(services *BlogServices) *BlogHandler {
	return &BlogHandler{services: services}
}
func (h *BlogHandler) GetBlogById() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		parsedId, err := strconv.Atoi(id)
		if err != nil {

		}

	})
}

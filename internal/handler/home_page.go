package handler

import (
	"lang/pkg/utils"
	"net/http"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"home": "this page is for front-end"}

	utils.ResponseServer(response, w, h.logger)
}

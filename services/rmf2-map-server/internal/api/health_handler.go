package api

import (
	"net/http"

	"rmf2map/server/internal/response"
)

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "rmf2 map backend is running"})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) GetOrganisation(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, map[string]string{"organisation": h.org(r)})
}

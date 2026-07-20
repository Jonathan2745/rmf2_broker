package api

import (
	"net/http"

	"rmf2map/server/internal/response"
	"rmf2map/server/internal/service"
)

func (h *Handler) GetMap(w http.ResponseWriter, r *http.Request) {
	lif, err := service.LoadMap(h.dataDir, h.org(r))
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, lif)
}

package api

import (
	"net/http"

	"rmf2map/server/internal/repository"
	"rmf2map/server/internal/response"
)

func (h *Handler) GetScene(w http.ResponseWriter, r *http.Request) {
	response.ServeGLB(w, repository.ScenePath(h.dataDir, h.org(r)), "")
}

func (h *Handler) GetModel(w http.ResponseWriter, r *http.Request) {
	model := r.PathValue("robotModel")
	fallback := repository.DefaultModelPath(h.dataDir, h.org(r))
	response.ServeGLB(w, repository.ModelPath(h.dataDir, h.org(r), model), fallback)
}

package api

import (
	"net/http"

	"rmf2map/server/internal/config"
	"rmf2map/server/internal/vda5050"
)

func RegisterRoutes(mux *http.ServeMux, cfg config.Config, poller *vda5050.Poller) {
	h := NewHandler(cfg.DataDir, cfg.DefaultOrg, poller)

	mux.HandleFunc("GET /organisation", h.GetOrganisation)
	mux.HandleFunc("GET /robots", h.GetRobots)
	mux.HandleFunc("GET /scene", h.GetScene)
	mux.HandleFunc("GET /map", h.GetMap)
	mux.HandleFunc("GET /path/{robotId}", h.GetPath)
	mux.HandleFunc("GET /models/{robotModel}", h.GetModel)
	mux.HandleFunc("GET /position/{robotId}", h.GetPosition)
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /{$}", h.Root)
}

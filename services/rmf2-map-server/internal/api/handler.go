package api

import (
	"net/http"

	"rmf2-map-server/internal/middleware"
	"rmf2-map-server/internal/vda5050"
)

type Handler struct {
	dataDir    string
	defaultOrg string
	poller     *vda5050.Poller // nil when VDA5050 integration is disabled
}

func NewHandler(dataDir, defaultOrg string, poller *vda5050.Poller) *Handler {
	return &Handler{dataDir: dataDir, defaultOrg: defaultOrg, poller: poller}
}

func (h *Handler) org(r *http.Request) string {
	if org := middleware.OrgFromContext(r.Context()); org != "" {
		return org
	}
	return h.defaultOrg
}

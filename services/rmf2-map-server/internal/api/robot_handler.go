package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"rmf2map/server/internal/response"
	"rmf2map/server/internal/service"
)

func (h *Handler) GetRobots(w http.ResponseWriter, r *http.Request) {
	robots, err := service.LoadRobots(h.dataDir, h.org(r))
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	type robotSummary struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Model string `json:"model"`
	}
	out := make([]robotSummary, len(robots))
	for i, rb := range robots {
		out[i] = robotSummary{ID: rb.ID, Name: rb.Name, Model: rb.Model}
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{"robots": out})
}

func (h *Handler) GetPath(w http.ResponseWriter, r *http.Request) {
	robotID, err := strconv.Atoi(r.PathValue("robotId"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "robotId must be an integer")
		return
	}
	nodes, edges, loop, err := service.LoadPath(h.dataDir, h.org(r), robotID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, fmt.Sprintf("path not found: %v", err))
		return
	}
	type pathEdge struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	edgeOut := make([]pathEdge, len(edges))
	for i, e := range edges {
		edgeOut[i] = pathEdge{From: e[0], To: e[1]}
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{
		"nodes": nodes,
		"edges": edgeOut,
		"loop":  loop,
	})
}

// positionJSON converts robotics coordinates (Y = ground plane) to Three.js/WebGL
// coordinates (Y = up axis, floor is X-Z plane) by mapping robotics-Y → z and setting y = 0.
func positionJSON(id int, x, y, theta float64, state string) map[string]any {
	out := map[string]any{
		"id":    id,
		"x":     x,
		"y":     y,
		"z":     10,
		"theta": theta,
	}
	if state != "" {
		out["state"] = state
	}
	return out
}

func (h *Handler) GetPosition(w http.ResponseWriter, r *http.Request) {
	robotID, err := strconv.Atoi(r.PathValue("robotId"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "robotId must be an integer")
		return
	}

	// Try live VDA5050 position first when the poller is active.
	if h.poller != nil {
		robots, err := service.LoadRobots(h.dataDir, h.org(r))
		if err == nil {
			for _, rb := range robots {
				if rb.ID == robotID && rb.VDA5050ID != "" {
					if pos, ok := h.poller.Position(rb.VDA5050ID); ok {
						response.WriteJSON(w, http.StatusOK, positionJSON(robotID, pos.X, pos.Y, pos.Theta, pos.State))
						return
					}
				}
			}
		}
	}

	if pos, ok, err := service.LoadSimulatedPosition(h.dataDir, h.org(r), robotID, time.Now()); err == nil && ok {
		response.WriteJSON(w, http.StatusOK, positionJSON(robotID, pos.X, pos.Y, pos.Theta, pos.State))
		return
	}

	// Fall back to static positions.json.
	pos, err := service.LoadPosition(h.dataDir, h.org(r), robotID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, fmt.Sprintf("position not found: %v", err))
		return
	}
	response.WriteJSON(w, http.StatusOK, positionJSON(robotID, pos.X, pos.Y, pos.Theta, pos.State))
}

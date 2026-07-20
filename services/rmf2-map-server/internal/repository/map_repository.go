package repository

import "path/filepath"

type Node struct {
	NodeID string  `json:"node_id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Theta  float64 `json:"theta"`
}

type Edge struct {
	EdgeID        string `json:"edge_id"`
	StartNodeID   string `json:"start_node_id"`
	EndNodeID     string `json:"end_node_id"`
	Bidirectional bool   `json:"bidirectional"`
}

type LIFFile struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

func MapPath(dataDir, org string) string { return filepath.Join(dataDir, org, "map.lif.json") }

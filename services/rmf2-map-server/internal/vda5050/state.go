package vda5050

// AGVState is a single robot's live state as returned by the VDA5050 master GET /state.
type AGVState struct {
	RobotID    string  `json:"robot_id"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Theta      float64 `json:"theta"`
	Driving    bool    `json:"driving"`
	HasPose    bool    `json:"has_pose"`
	LastNodeID string  `json:"last_node_id"`
}

// StateResponse is the full response body of GET /state.
type StateResponse struct {
	Type string     `json:"type"`
	AGVs []AGVState `json:"agvs"`
}

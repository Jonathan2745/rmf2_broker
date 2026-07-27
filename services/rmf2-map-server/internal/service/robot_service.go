package service

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"rmf2-map-server/internal/repository"
)

const defaultSpeedMPS = 1.0

func LoadAPIKeys(dataDir string) (repository.APIKeys, error) {
	data, err := os.ReadFile(repository.APIKeysPath(dataDir))
	if err != nil {
		return nil, fmt.Errorf("read api-keys.json: %w", err)
	}
	var keys repository.APIKeys
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, fmt.Errorf("parse api-keys.json: %w", err)
	}
	return keys, nil
}

func LoadRobots(dataDir, org string) ([]repository.Robot, error) {
	var rf repository.RobotsFile
	if err := repository.ReadJSONFile(repository.RobotsPath(dataDir, org), &rf); err != nil {
		return nil, err
	}
	return rf.Robots, nil
}

func LoadMap(dataDir, org string) (repository.LIFFile, error) {
	var lif repository.LIFFile
	if err := repository.ReadJSONFile(repository.MapPath(dataDir, org), &lif); err != nil {
		return repository.LIFFile{}, err
	}
	return lif, nil
}

// LoadPath returns the robot's waypoints, loop flag, and sequential edge pairs
// e.g. [["WP0","WP1"],["WP1","WP2"],...].
func LoadPath(dataDir, org string, robotID int) ([]repository.Waypoint, [][2]string, bool, error) {
	robots, err := LoadRobots(dataDir, org)
	if err != nil {
		return nil, nil, false, err
	}
	for _, r := range robots {
		if r.ID != robotID {
			continue
		}
		wps := r.Path
		edgeCount := 0
		if len(wps) > 1 {
			edgeCount = len(wps) - 1
			if r.Loop {
				edgeCount++
			}
		}
		edges := make([][2]string, 0, edgeCount)
		for i := 1; i < len(wps); i++ {
			edges = append(edges, [2]string{wps[i-1].ID, wps[i].ID})
		}
		if r.Loop && len(wps) > 1 {
			edges = append(edges, [2]string{wps[len(wps)-1].ID, wps[0].ID})
		}
		return wps, edges, r.Loop, nil
	}
	return nil, nil, false, fmt.Errorf("robot %d not found", robotID)
}

func LoadSimulatedPosition(dataDir, org string, robotID int, now time.Time) (repository.Position, bool, error) {
	robots, err := LoadRobots(dataDir, org)
	if err != nil {
		return repository.Position{}, false, err
	}
	for _, r := range robots {
		if r.ID == robotID {
			pos, ok := simulatedPosition(r, now)
			return pos, ok, nil
		}
	}
	return repository.Position{}, false, nil
}

type routeSegment struct {
	from   repository.Waypoint
	to     repository.Waypoint
	length float64
}

func simulatedPosition(r repository.Robot, now time.Time) (repository.Position, bool) {
	if len(r.Path) == 0 {
		return repository.Position{}, false
	}
	if len(r.Path) == 1 {
		wp := r.Path[0]
		return repository.Position{X: wp.X, Y: wp.Y, Theta: wp.Theta, State: "stopped"}, true
	}

	segments := make([]routeSegment, 0, len(r.Path))
	totalLength := 0.0
	for i := 1; i < len(r.Path); i++ {
		seg := newRouteSegment(r.Path[i-1], r.Path[i])
		if seg.length == 0 {
			continue
		}
		segments = append(segments, seg)
		totalLength += seg.length
	}
	if r.Loop {
		seg := newRouteSegment(r.Path[len(r.Path)-1], r.Path[0])
		if seg.length > 0 {
			segments = append(segments, seg)
			totalLength += seg.length
		}
	}
	if totalLength == 0 {
		wp := r.Path[0]
		return repository.Position{X: wp.X, Y: wp.Y, Theta: wp.Theta, State: "stopped"}, true
	}

	speed := r.SpeedMPS
	if speed <= 0 {
		speed = defaultSpeedMPS
	}
	distance := now.Sub(time.Unix(0, 0)).Seconds() * speed
	state := "driving"
	if r.Loop {
		distance = math.Mod(distance, totalLength)
	} else if distance >= totalLength {
		last := r.Path[len(r.Path)-1]
		return repository.Position{X: last.X, Y: last.Y, Theta: last.Theta, State: "stopped"}, true
	}

	for _, seg := range segments {
		if distance > seg.length {
			distance -= seg.length
			continue
		}
		ratio := distance / seg.length
		return repository.Position{
			X:     seg.from.X + (seg.to.X-seg.from.X)*ratio,
			Y:     seg.from.Y + (seg.to.Y-seg.from.Y)*ratio,
			Theta: math.Atan2(seg.to.Y-seg.from.Y, seg.to.X-seg.from.X),
			State: state,
		}, true
	}

	last := segments[len(segments)-1].to
	return repository.Position{X: last.X, Y: last.Y, Theta: last.Theta, State: state}, true
}

func newRouteSegment(from, to repository.Waypoint) routeSegment {
	return routeSegment{
		from:   from,
		to:     to,
		length: math.Hypot(to.X-from.X, to.Y-from.Y),
	}
}

func LoadPosition(dataDir, org string, robotID int) (repository.Position, error) {
	data, err := os.ReadFile(repository.PositionsPath(dataDir, org))
	if err != nil {
		return repository.Position{}, fmt.Errorf("read positions.json: %w", err)
	}
	var raw map[string]repository.Position
	if err := json.Unmarshal(data, &raw); err != nil {
		return repository.Position{}, fmt.Errorf("parse positions.json: %w", err)
	}
	key := fmt.Sprintf("%d", robotID)
	pos, ok := raw[key]
	if !ok {
		return repository.Position{}, fmt.Errorf("robot %d not found in positions", robotID)
	}
	return pos, nil
}

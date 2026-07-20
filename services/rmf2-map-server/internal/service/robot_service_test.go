package service

import (
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadPathIncludesLoopClosingEdge(t *testing.T) {
	dataDir := writeRobotsFile(t)

	_, edges, loop, err := LoadPath(dataDir, "test-org", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !loop {
		t.Fatal("expected loop flag")
	}

	want := [][2]string{
		{"A", "B"},
		{"B", "C"},
		{"C", "D"},
		{"D", "A"},
	}
	if len(edges) != len(want) {
		t.Fatalf("expected %d edges, got %d", len(want), len(edges))
	}
	for i := range want {
		if edges[i] != want[i] {
			t.Fatalf("edge %d: expected %#v, got %#v", i, want[i], edges[i])
		}
	}
}

func TestLoadSimulatedPositionLoopsAroundRoute(t *testing.T) {
	dataDir := writeRobotsFile(t)

	pos, ok, err := LoadSimulatedPosition(dataDir, "test-org", 1, time.Unix(5, 0))
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected simulated position")
	}

	assertFloat(t, pos.X, 10)
	assertFloat(t, pos.Y, 0)
	assertFloat(t, pos.Theta, 0)
	if pos.State != "driving" {
		t.Fatalf("expected driving state, got %q", pos.State)
	}

	pos, ok, err = LoadSimulatedPosition(dataDir, "test-org", 1, time.Unix(20, 0))
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected simulated position")
	}
	assertFloat(t, pos.X, 0)
	assertFloat(t, pos.Y, 0)
}

func writeRobotsFile(t *testing.T) string {
	t.Helper()

	dataDir := t.TempDir()
	orgDir := filepath.Join(dataDir, "test-org")
	if err := os.MkdirAll(orgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	const robotsJSON = `{
  "robots": [
    {
      "id": 1,
      "name": "Robot 01",
      "model": "AMR-Type1",
      "speed_mps": 2,
      "loop": true,
      "path": [
        {"id": "A", "x": 0, "y": 0, "theta": 0, "label": "A"},
        {"id": "B", "x": 10, "y": 0, "theta": 0, "label": "B"},
        {"id": "C", "x": 10, "y": 10, "theta": 0, "label": "C"},
        {"id": "D", "x": 0, "y": 10, "theta": 0, "label": "D"}
      ]
    }
  ]
}`
	if err := os.WriteFile(filepath.Join(orgDir, "robots.json"), []byte(robotsJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	return dataDir
}

func assertFloat(t *testing.T, got, want float64) {
	t.Helper()

	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

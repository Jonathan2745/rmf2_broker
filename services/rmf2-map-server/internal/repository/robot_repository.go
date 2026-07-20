package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Robot struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Model     string     `json:"model"`
	VDA5050ID string     `json:"vda5050_id,omitempty"`
	SpeedMPS  float64    `json:"speed_mps,omitempty"`
	Loop      bool       `json:"loop,omitempty"`
	Path      []Waypoint `json:"path"`
}

type Waypoint struct {
	ID    string  `json:"id"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Theta float64 `json:"theta"`
	Label string  `json:"label"`
}

type RobotsFile struct {
	Robots []Robot `json:"robots"`
}

type Position struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Theta float64 `json:"theta"`
	State string  `json:"state,omitempty"`
}

type APIKeys map[string]string // bearer token → org name

func OrgDir(dataDir, org string) string        { return filepath.Join(dataDir, org) }
func RobotsPath(dataDir, org string) string    { return filepath.Join(dataDir, org, "robots.json") }
func PositionsPath(dataDir, org string) string { return filepath.Join(dataDir, org, "positions.json") }
func APIKeysPath(dataDir string) string        { return filepath.Join(dataDir, "api-keys.json") }

func ReadJSONFile(path string, dest any) error {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("open %s: %w", filepath.Base(path), err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(dest); err != nil {
		return fmt.Errorf("invalid JSON in %s: %w", filepath.Base(path), err)
	}
	return nil
}

func WriteJSONFile(path string, src any) error {
	b, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return fmt.Errorf("write tmp: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}

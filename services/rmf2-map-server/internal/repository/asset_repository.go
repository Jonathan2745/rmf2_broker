package repository

import "path/filepath"

func ScenePath(dataDir, org string) string { return filepath.Join(dataDir, org, "scene.draco.glb") }
func ModelPath(dataDir, org, model string) string {
	return filepath.Join(dataDir, org, "models", model+".glb")
}
func DefaultModelPath(dataDir, org string) string {
	return filepath.Join(dataDir, org, "models", "robot.glb")
}

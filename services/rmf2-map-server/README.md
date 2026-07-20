# RMF2 Map Backend

A multi-tenant Go HTTP server that serves robot fleet data — map graphs, robot metadata, planned paths, live positions, and 3D assets — scoped per organisation via API key authentication.

---

## Requirements

- Go 1.22+

---

## Running

```bash
go run ./cmd/server
```

**Flags:**

| Flag | Default | Description |
|---|---|---|
| `-addr` | `0.0.0.0:8008` | Listen address |
| `-data-dir` | `./data` | Path to data directory |
| `-default-org` | `ros-industrial` | Fallback organisation |

---

## Authentication

All routes except `/health` require a Bearer token in the `Authorization` header.

```
Authorization: Bearer <api-key>
```

Keys are mapped to organisations in `data/api-keys.json`:

```json
{
  "key-ros-industrial-secret": "ros-industrial",
  "key-acme-robotics-secret": "acme-robotics"
}
```

The resolved organisation determines which data directory is used. A caller cannot access another org's data.

---

## API Reference

### `GET /organisation`
Returns the organisation the caller's API key belongs to.

```json
{ "organisation": "ros-industrial" }
```

---

### `GET /robots`
Returns the list of robots for the caller's organisation.

```json
{
  "robots": [
    { "id": 1, "name": "Robot 01", "model": "AMR-Type1" },
    { "id": 2, "name": "Robot 02", "model": "AGV-Type2" },
    { "id": 3, "name": "Robot 03", "model": "AMR-Type1" }
  ]
}
```

---

### `GET /scene`
Streams the organisation's scene as a binary GLB file (`Content-Type: model/gltf-binary`).  
Expects `data/{org}/scene.draco.glb` on disk. Returns `404` if not found.

---

### `GET /map`
Returns the full factory floor graph — all nodes and edges from the LIF file.

```json
{
  "nodes": [
    { "node_id": "P0", "x": -30.8, "y": -49.3, "theta": 0 },
    ...
  ],
  "edges": [
    { "edge_id": "E1", "start_node_id": "P0", "end_node_id": "P1", "bidirectional": true },
    ...
  ]
}
```

---

### `GET /path/{robotId}`
Returns the planned path for a robot as an ordered list of waypoints and the sequential edges connecting them.
Looped dummy routes include a closing edge from the final waypoint back to the first waypoint.

```json
{
  "loop": true,
  "nodes": [
    { "id": "WP0", "x": -10.0, "y": 25.5, "theta": 0.0, "label": "Start" },
    { "id": "WP1", "x": -0.5,  "y": 16.0, "theta": 0.0, "label": "WP 1" }
  ],
  "edges": [
    { "from": "WP0", "to": "WP1" }
  ]
}
```

Returns `404` if the robot ID does not exist.

---

### `GET /models/{robotModel}`
Streams the GLB model file for a given robot model string (`Content-Type: model/gltf-binary`).  
Looks for `data/{org}/models/{robotModel}.glb`. Falls back to `data/{org}/models/robot.glb` if the specific model file is not found.

---

### `GET /position/{robotId}`
Returns the current position, orientation, and state of a robot. Intended to be polled every 500ms.
When VDA5050 is disabled, dummy positions are simulated from each robot's route, speed, and loop settings in `robots.json`.

```json
{ "id": 1, "x": 12.0, "y": 0.0, "z": 1.5, "theta": 0.0, "state": "stopped" }
```

Returns `404` if the robot ID does not exist.

---

### `GET /health`
Health check — no authentication required.

```json
{ "status": "ok" }
```

---

## Data Directory Structure

Each organisation has its own subdirectory under `data/`:

```
data/
  api-keys.json              # bearer token → org name mapping
  {org}/
    robots.json              # robot list with model strings, speed, loop flag, and path waypoints
    map.lif.json             # LIF graph (nodes + edges)
    positions.json           # static position fallback per robot (keyed by robot id)
    scene.draco.glb          # scene model (user-provided)
    models/
      robot.glb              # default fallback model (user-provided)
      {robotModel}.glb       # per-model GLB files (user-provided)
```

### `robots.json` format

```json
{
  "robots": [
    {
      "id": 1,
      "name": "Robot 01",
      "model": "AMR-Type1",
      "speed_mps": 1.2,
      "loop": true,
      "path": [
        { "id": "R1-A", "x": -10.0, "y": 16.0, "theta": 0.0, "label": "Square A" },
        { "id": "R1-B", "x": 12.0,  "y": 16.0, "theta": 0.0, "label": "Square B" },
        { "id": "R1-C", "x": 12.0,  "y": 4.0,  "theta": 0.0, "label": "Square C" },
        { "id": "R1-D", "x": -10.0, "y": 4.0,  "theta": 0.0, "label": "Square D" }
      ]
    }
  ]
}
```

`speed_mps` controls dummy simulation speed. `loop: true` makes the robot travel from the final waypoint back to the first waypoint and continue indefinitely.

### `positions.json` format

```json
{
  "1": { "x": 12.0, "y": 4.0, "theta": 0.0, "state": "stopped" },
  "2": { "x": 26.0, "y": 4.0, "theta": 0.0, "state": "stopped" }
}
```

---

## Adding a New Organisation

1. Create the directory `data/{org}/` with the file structure above.
2. Add an entry to `data/api-keys.json`:
   ```json
   "your-secret-key": "your-org-name"
   ```
3. Restart the server — keys are loaded once at startup.

---

## Project Structure

```
cmd/server/main.go                # entry point — flags, key loading, wiring
internal/
  api/
    handler.go                    # HTTP handler methods
    routes.go                     # URL → handler mapping
  config/config.go                # Config struct
  middleware/auth.go              # CORS + Bearer token auth
  repository/robot_repository.go  # data types + file path helpers + JSON I/O
  response/json.go                # WriteJSON, WriteError, ServeGLB helpers
  service/robot_service.go        # business logic — load robots, map, path, position
data/
  api-keys.json
  ros-industrial/                 # example organisation
```

package vda5050

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"rmf2map/server/internal/repository"
)

const pollInterval = 200 * time.Millisecond

// Poller fetches live AGV state from the VDA5050 master on a fixed interval
// and caches the latest position for each robot.
type Poller struct {
	base   string
	client *http.Client
	mu     sync.RWMutex
	cache  map[string]repository.Position // vda5050_id → position
}

func NewPoller(base string) *Poller {
	return &Poller{
		base:   base,
		client: &http.Client{Timeout: 5 * time.Second},
		cache:  make(map[string]repository.Position),
	}
}

// Start launches the background polling goroutine. It stops when ctx is cancelled.
func (p *Poller) Start(ctx context.Context) {
	if err := p.fetch(); err != nil {
		log.Printf("vda5050 initial fetch: %v", err)
	}

	go func() {
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := p.fetch(); err != nil {
					log.Printf("vda5050 poll: %v", err)
				}
			}
		}
	}()
}

func (p *Poller) fetch() error {
	resp, err := p.client.Get(p.base + "/state")
	if err != nil {
		return fmt.Errorf("GET /state: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET /state: status %d", resp.StatusCode)
	}

	var state StateResponse
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	p.mu.Lock()
	for _, agv := range state.AGVs {
		if !agv.HasPose {
			continue
		}
		p.cache[agv.RobotID] = repository.Position{
			X:     agv.X,
			Y:     agv.Y,
			Theta: agv.Theta,
			State: robotState(agv.Driving),
		}
	}
	p.mu.Unlock()
	return nil
}

func robotState(driving bool) string {
	if driving {
		return "driving"
	}
	return "stopped"
}

// Position returns the latest cached position for the given VDA5050 robot ID.
// Returns false if no data is available yet.
func (p *Poller) Position(vda5050ID string) (repository.Position, bool) {
	p.mu.RLock()
	pos, ok := p.cache[vda5050ID]
	p.mu.RUnlock()
	return pos, ok
}

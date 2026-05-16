package mcp

import (
	"context"
	"log"
	"time"

	"StreamCore/internal/pkg/db/ai"
)

// StartPeriodicSync periodically syncs all active MCP servers. Blocks until ctx is cancelled.
// Caller should invoke in a goroutine: go StartPeriodicSync(ctx, registry, db).
func StartPeriodicSync(ctx context.Context, registry *ToolRegistry, db ai.AIDatabase) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Printf("[ai] mcp sync started")
	for {
		select {
		case <-ctx.Done():
			log.Printf("[ai] mcp sync stopped")
			return
		case <-ticker.C:
			syncAllServers(ctx, registry, db)
		}
	}
}

func syncAllServers(ctx context.Context, registry *ToolRegistry, db ai.AIDatabase) {
	servers, err := db.ListServers(ctx)
	if err != nil {
		log.Printf("[ai] mcp sync: list servers failed: %v", err)
		return
	}

	for _, s := range servers {
		if s.Status != 1 {
			continue
		}
		interval := time.Duration(s.SyncIntervalSec) * time.Second
		if interval <= 0 {
			interval = 5 * time.Minute
		}
		if s.LastSyncedAt != nil && time.Since(*s.LastSyncedAt) < interval {
			continue
		}

		log.Printf("[ai] mcp sync: syncing server %s (%s)", s.ServerName, s.ServerURL)
		if err := registry.SyncServer(ctx, s); err != nil {
			log.Printf("[ai] mcp sync: server %s failed: %v", s.ServerName, err)
		}
	}
}

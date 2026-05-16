package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"StreamCore/internal/pkg/ai/llm"
	"StreamCore/internal/pkg/db/ai"
	"StreamCore/internal/pkg/db/model"
)

// ToolRegistry manages MCP tool discovery and execution.
type ToolRegistry struct {
	db     ai.AIDatabase
	client Client
}

// NewToolRegistry creates a ToolRegistry.
func NewToolRegistry(db ai.AIDatabase, client Client) *ToolRegistry {
	return &ToolRegistry{db: db, client: client}
}

// SyncServer discovers tools from an MCP server and upserts them into the DB.
func (r *ToolRegistry) SyncServer(ctx context.Context, server *model.MCPServerModel) error {
	tools, err := r.client.DiscoverTools(ctx, server.ServerURL)
	if err != nil {
		return fmt.Errorf("sync server %s: %w", server.ServerName, err)
	}

	// Clear old tools for this server
	if err := r.db.DeleteToolsByServer(ctx, server.ID); err != nil {
		return fmt.Errorf("sync server %s: delete old tools: %w", server.ServerName, err)
	}

	// Upsert new tools
	for _, t := range tools {
		now := time.Now()
		tool := &model.MCPToolModel{
			ToolID:      server.ServerName + "/" + t.Name,
			ToolName:    t.Name,
			Description: t.Description,
			InputSchema: string(t.InputSchema),
			ServerID:    server.ID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := r.db.UpsertTool(ctx, tool); err != nil {
			return fmt.Errorf("sync server %s: upsert tool %s: %w", server.ServerName, t.Name, err)
		}
	}

	// Update last synced time
	now := time.Now()
	server.LastSyncedAt = &now
	if err := r.db.UpdateServer(ctx, server); err != nil {
		return fmt.Errorf("sync server %s: update last_synced_at: %w", server.ServerName, err)
	}

	return nil
}

// ExecuteTool runs a specific tool by loading it from DB and calling the MCP server.
func (r *ToolRegistry) ExecuteTool(ctx context.Context, toolID string, args json.RawMessage) (string, error) {
	tools, err := r.db.ListToolsByIDs(ctx, []string{toolID})
	if err != nil {
		return "", fmt.Errorf("execute tool %s: %w", toolID, err)
	}
	if len(tools) == 0 {
		return "", fmt.Errorf("execute tool %s: not found", toolID)
	}
	tool := tools[0]

	server, err := r.db.GetServer(ctx, tool.ServerID)
	if err != nil {
		return "", fmt.Errorf("execute tool %s: get server: %w", toolID, err)
	}

	return r.client.CallTool(ctx, server.ServerURL, tool.ToolName, args)
}

// FilterTools returns LLM-compatible tool definitions for the given tool IDs.
func (r *ToolRegistry) FilterTools(ctx context.Context, toolIDs []string) ([]llm.ToolDef, error) {
	tools, err := r.db.ListToolsByIDs(ctx, toolIDs)
	if err != nil {
		return nil, fmt.Errorf("filter tools: %w", err)
	}

	result := make([]llm.ToolDef, 0, len(tools))
	for _, t := range tools {
		result = append(result, llm.ToolDef{
			Name:        t.ToolID,
			Description: t.Description,
			InputSchema: json.RawMessage(t.InputSchema),
		})
	}
	return result, nil
}

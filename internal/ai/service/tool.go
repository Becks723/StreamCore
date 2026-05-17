package service

import (
	"fmt"

	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/pack"
	kitexai "StreamCore/kitex_gen/ai"
	kitexcommon "StreamCore/kitex_gen/common"
)

func RegisterMCPServer(s *AIService, req *kitexai.RegisterMCPServerReq) (*kitexai.RegisterMCPServerResp, error) {
	syncInterval := int(300)
	if req.IsSetSyncIntervalSec() {
		syncInterval = int(req.GetSyncIntervalSec())
	}

	server := &model.MCPServerModel{
		ServerName:      req.ServerName,
		ServerURL:       req.ServerUrl,
		SyncIntervalSec: syncInterval,
		Status:          1,
	}
	if req.IsSetAuthToken() {
		server.AuthToken = req.GetAuthToken()
	}

	if err := s.db.CreateServer(s.ctx, server); err != nil {
		return nil, fmt.Errorf("RegisterMCPServer: %w", err)
	}

	// Discover tools from the MCP server and store them
	if err := s.toolReg.SyncServer(s.ctx, server); err != nil {
		return nil, fmt.Errorf("RegisterMCPServer: sync tools: %w", err)
	}

	// Reload tools to return in response
	tools, _ := s.db.ListToolsByIDs(s.ctx, nil) // refresh after discovery — not needed for resp, skip
	_ = tools

	resp := new(kitexai.RegisterMCPServerResp)
	resp.Base = pack.BuildSuccessResp()
	resp.Server = pack.MCPServerInfo(server)
	return resp, nil
}

func RefreshMCPServer(s *AIService, req *kitexai.RefreshMCPServerReq) (*kitexai.RefreshMCPServerResp, error) {
	serverID := uint(req.ServerId)
	server, err := s.db.GetServer(s.ctx, serverID)
	if err != nil {
		return nil, fmt.Errorf("RefreshMCPServer: get server failed: %w", err)
	}

	// Sync tools from the MCP server
	if err := s.toolReg.SyncServer(s.ctx, server); err != nil {
		return nil, fmt.Errorf("RefreshMCPServer: sync tools: %w", err)
	}

	// Reload tools to return in response
	updatedTools, err := s.db.ListTools(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("RefreshMCPServer: list updated tools: %w", err)
	}
	serverName := server.ServerName
	items := make([]*kitexcommon.MCPToolInfo, 0, len(updatedTools))
	for _, t := range updatedTools {
		if t.ServerID == serverID {
			items = append(items, pack.MCPToolInfo(t, serverName))
		}
	}

	resp := new(kitexai.RefreshMCPServerResp)
	resp.Base = pack.BuildSuccessResp()
	resp.Tools = items
	return resp, nil
}

func ListMCPServers(s *AIService, req *kitexai.ListMCPServersReq) (*kitexai.ListMCPServersResp, error) {
	servers, err := s.db.ListServers(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("ListMCPServers: %w", err)
	}

	items := make([]*kitexcommon.MCPServerInfo, len(servers))
	for i, srv := range servers {
		items[i] = pack.MCPServerInfo(srv)
	}

	resp := new(kitexai.ListMCPServersResp)
	resp.Base = pack.BuildSuccessResp()
	resp.Data = &kitexai.MCPServerListData{Servers: items}
	return resp, nil
}

func DeleteMCPServer(s *AIService, req *kitexai.DeleteMCPServerReq) (*kitexai.DeleteMCPServerResp, error) {
	serverID := uint(req.ServerId)
	if err := s.db.DeleteServer(s.ctx, serverID); err != nil {
		return nil, fmt.Errorf("DeleteMCPServer: %w", err)
	}

	resp := new(kitexai.DeleteMCPServerResp)
	resp.Base = pack.BuildSuccessResp()
	return resp, nil
}

func ListTools(s *AIService, req *kitexai.ListToolsReq) (*kitexai.ListToolsResp, error) {
	tools, err := s.db.ListTools(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("ListTools: %w", err)
	}

	serverMap := make(map[uint]string)
	servers, _ := s.db.ListServers(s.ctx)
	for _, srv := range servers {
		serverMap[srv.ID] = srv.ServerName
	}

	items := make([]*kitexcommon.MCPToolInfo, len(tools))
	for i, t := range tools {
		items[i] = pack.MCPToolInfo(t, serverMap[t.ServerID])
	}

	resp := new(kitexai.ListToolsResp)
	resp.Base = pack.BuildSuccessResp()
	resp.Data = &kitexai.ToolListData{Tools: items}
	return resp, nil
}

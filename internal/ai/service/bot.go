package service

import (
	"fmt"

	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/pack"
	kitexai "StreamCore/kitex_gen/ai"
	kitexcommon "StreamCore/kitex_gen/common"
	"StreamCore/pkg/util"
	"gorm.io/gorm"
)

func CreateBot(s *AIService, req *kitexai.CreateBotReq) (*kitexai.CreateBotResp, error) {
	// Create a user record for the bot
	bot := &model.UserModel{
		Username: req.GetBotName(),
		IsBot:    true,
		BotConfig: pack.BotConfigToJSON(&pack.BotConfigData{
			SystemPrompt: req.GetSystemPrompt(),
			ModelName:    req.GetModelName(),
			TriggerMode:  req.GetTriggerMode(),
			ToolIDs:      req.GetToolIds(),
		}),
	}

	if err := s.db.CreateBotUser(s.ctx, bot); err != nil {
		return nil, fmt.Errorf("CreateBot: create user failed: %w", err)
	}

	info := pack.BotInfo(bot)
	resp := new(kitexai.CreateBotResp)
	resp.Base = pack.BuildSuccessResp()
	resp.Bot = info
	return resp, nil
}

func UpdateBot(s *AIService, req *kitexai.UpdateBotReq) (*kitexai.UpdateBotResp, error) {
	botID, err := util.ParseUint(req.GetBotId())
	if err != nil {
		return nil, fmt.Errorf("UpdateBot: bad bot_id: %w", err)
	}

	bot, err := s.db.GetBotUser(s.ctx, botID)
	if err != nil {
		return nil, fmt.Errorf("UpdateBot: get bot user failed: %w", err)
	}

	cfg := pack.ParseBotConfig(bot.BotConfig)

	if req.IsSetBotName() {
		bot.Username = req.GetBotName()
	}
	if req.IsSetSystemPrompt() {
		cfg.SystemPrompt = req.GetSystemPrompt()
	}
	if req.IsSetModelName() {
		cfg.ModelName = req.GetModelName()
	}
	if req.IsSetTriggerMode() {
		cfg.TriggerMode = req.GetTriggerMode()
	}
	if req.IsSetToolIds() {
		cfg.ToolIDs = req.GetToolIds()
	}

	bot.BotConfig = pack.BotConfigToJSON(cfg)

	if err := s.db.UpdateBotUser(s.ctx, bot); err != nil {
		return nil, fmt.Errorf("UpdateBot: update user failed: %w", err)
	}

	info := pack.BotInfo(bot)
	resp := new(kitexai.UpdateBotResp)
	resp.Base = pack.BuildSuccessResp()
	resp.Bot = info
	return resp, nil
}

func ListBots(s *AIService, req *kitexai.ListBotsReq) (*kitexai.ListBotsResp, error) {
	page := 1
	pageSize := 20
	if req.IsSetPage() {
		page = int(req.GetPage())
	}
	if req.IsSetPageSize() {
		pageSize = int(req.GetPageSize())
	}

	bots, total, err := s.db.ListBotUsers(s.ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("ListBots: %w", err)
	}

	items := make([]*kitexcommon.BotInfo, len(bots))
	for i, b := range bots {
		items[i] = pack.BotInfo(b)
	}

	resp := new(kitexai.ListBotsResp)
	resp.Base = pack.BuildSuccessResp()
	resp.Bots = items
	t := int32(total)
	resp.Total = &t
	return resp, nil
}

func GetBot(s *AIService, req *kitexai.GetBotReq) (*kitexai.GetBotResp, error) {
	botID, err := util.ParseUint(req.BotId)
	if err != nil {
		return nil, fmt.Errorf("GetBot: bad bot_id: %w", err)
	}
	bot, err := s.db.GetBotUser(s.ctx, botID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("GetBot: bot not found")
		}
		return nil, fmt.Errorf("GetBot: %w", err)
	}

	resp := new(kitexai.GetBotResp)
	resp.Base = pack.BuildSuccessResp()
	resp.Bot = pack.BotInfo(bot)
	return resp, nil
}

func DeleteBot(s *AIService, req *kitexai.DeleteBotReq) (*kitexai.DeleteBotResp, error) {
	botID, err := util.ParseUint(req.BotId)
	if err != nil {
		return nil, fmt.Errorf("DeleteBot: bad bot_id: %w", err)
	}

	if err := s.db.DeleteBotUser(s.ctx, botID); err != nil {
		return nil, fmt.Errorf("DeleteBot: delete user failed: %w", err)
	}

	resp := new(kitexai.DeleteBotResp)
	resp.Base = pack.BuildSuccessResp()
	return resp, nil
}

// Bot-Group binding methods

func AddBotToGroup(s *AIService, req *kitexai.AddBotToGroupReq) (*kitexai.AddBotToGroupResp, error) {
	botID, err := util.ParseUint(req.BotId)
	if err != nil {
		return nil, fmt.Errorf("AddBotToGroup: bad bot_id: %w", err)
	}
	groupID, err := util.ParseUint(req.GroupId)
	if err != nil {
		return nil, fmt.Errorf("AddBotToGroup: bad group_id: %w", err)
	}

	if err := s.db.AddBotToGroup(s.ctx, botID, groupID); err != nil {
		return nil, fmt.Errorf("AddBotToGroup: %w", err)
	}

	resp := new(kitexai.AddBotToGroupResp)
	resp.Base = pack.BuildSuccessResp()
	return resp, nil
}

func RemoveBotFromGroup(s *AIService, req *kitexai.RemoveBotFromGroupReq) (*kitexai.RemoveBotFromGroupResp, error) {
	botID, err := util.ParseUint(req.BotId)
	if err != nil {
		return nil, fmt.Errorf("RemoveBotFromGroup: bad bot_id: %w", err)
	}
	groupID, err := util.ParseUint(req.GroupId)
	if err != nil {
		return nil, fmt.Errorf("RemoveBotFromGroup: bad group_id: %w", err)
	}

	if err := s.db.RemoveBotFromGroup(s.ctx, botID, groupID); err != nil {
		return nil, fmt.Errorf("RemoveBotFromGroup: %w", err)
	}

	resp := new(kitexai.RemoveBotFromGroupResp)
	resp.Base = pack.BuildSuccessResp()
	return resp, nil
}

func ListGroupBots(s *AIService, req *kitexai.ListGroupBotsReq) (*kitexai.ListGroupBotsResp, error) {
	groupID, err := util.ParseUint(req.GroupId)
	if err != nil {
		return nil, fmt.Errorf("ListGroupBots: bad group_id: %w", err)
	}

	bots, err := s.db.ListGroupBots(s.ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("ListGroupBots: %w", err)
	}

	items := make([]*kitexcommon.BotInfo, len(bots))
	for i, b := range bots {
		items[i] = pack.BotInfo(b)
	}

	resp := new(kitexai.ListGroupBotsResp)
	resp.Base = pack.BuildSuccessResp()
	resp.Bots = items
	return resp, nil
}

func ListBotGroups(s *AIService, req *kitexai.ListBotGroupsReq) (*kitexai.ListBotGroupsResp, error) {
	botID, err := util.ParseUint(req.BotId)
	if err != nil {
		return nil, fmt.Errorf("ListBotGroups: bad bot_id: %w", err)
	}

	groupIDs, err := s.db.ListBotGroups(s.ctx, botID)
	if err != nil {
		return nil, fmt.Errorf("ListBotGroups: %w", err)
	}

	groupIDStrs := make([]string, len(groupIDs))
	for i, gid := range groupIDs {
		groupIDStrs[i] = util.Uint2String(gid)
	}

	resp := new(kitexai.ListBotGroupsResp)
	resp.Base = pack.BuildSuccessResp()
	resp.GroupIds = groupIDStrs
	return resp, nil
}

package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"StreamCore/internal/pkg/pack"
	kitexai "StreamCore/kitex_gen/ai"
	"StreamCore/pkg/util"
)

const (
	proactiveHistoryLimit = 20
	proactivePendingTTL   = time.Hour
)

// ProactiveCheck 记录群里最新一条真人消息，并登记一次延迟检查。
// 后续真人消息会覆盖 pending msg id，旧 timer 醒来后会自然失效。
// 返回 scheduled=false 表示当前群没有可用的 proactive bot，或本条消息来自 bot。
func ProactiveCheck(s *AIService, req *kitexai.ProactiveCheckReq) (*kitexai.ProactiveCheckResp, error) {
	resp := new(kitexai.ProactiveCheckResp)
	resp.Base = pack.BuildSuccessResp()
	resp.Scheduled = util.BoolPtr(false)

	groupID, err := util.ParseUint(req.GroupId)
	if err != nil {
		return nil, fmt.Errorf("ProactiveCheck: bad group_id: %w", err)
	}
	fromUID, err := util.ParseUint(req.FromUid)
	if err != nil {
		return nil, fmt.Errorf("ProactiveCheck: bad from_uid: %w", err)
	}

	botID, cfg, ok, err := selectProactiveBot(s.ctx, s, groupID, fromUID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return resp, nil
	}

	quietDelay := time.Duration(cfg.QuietMinutes) * time.Minute
	if err := s.infra.Cache.AI.SetProactivePending(s.ctx, groupID, req.MsgId, quietDelay+proactivePendingTTL); err != nil {
		return nil, fmt.Errorf("ProactiveCheck: set pending key: %w", err)
	}

	go func() {
		ctx := context.Background()
		timer := time.NewTimer(quietDelay)
		defer timer.Stop()
		<-timer.C

		if err := runProactiveCheck(ctx, s, groupID, botID, req.MsgId, cfg); err != nil {
			log.Printf("[ai] ProactiveCheck failed group=%d bot=%d msg=%d: %v", groupID, botID, req.MsgId, err)
		}
	}()

	resp.Scheduled = util.BoolPtr(true)
	return resp, nil
}

// selectProactiveBot 从群绑定 bot 中选出第一个开启 proactive 的 bot。
// 如果触发消息来自任意 bot，直接返回 ok=false，避免 bot 消息继续触发主动回复。
func selectProactiveBot(ctx context.Context, s *AIService, groupID, fromUID uint) (uint, *pack.ProactiveConfig, bool, error) {
	bots, err := s.db.ListGroupBots(ctx, groupID)
	if err != nil {
		return 0, nil, false, fmt.Errorf("ProactiveCheck: list group bots: %w", err)
	}
	for _, bot := range bots {
		// Bot 自己发出的消息不能继续触发主动回复，避免自我续聊。
		if bot.ID == fromUID {
			return 0, nil, false, nil
		}
		cfg := pack.ParseBotConfig(bot.BotConfig)
		if cfg.Proactive.Enabled {
			return bot.ID, &cfg.Proactive, true, nil
		}
	}
	return 0, nil, false, nil
}

// runProactiveCheck 在 quietDelay 后执行真正的冷场判断。
// 它会先确认 pending msg id 仍然等于本次消息，确保期间没有新的真人消息插入。
// 只有在频控允许、生成了非空回复时，才会以 bot 身份写回群聊。
func runProactiveCheck(ctx context.Context, s *AIService, groupID, botID uint, msgID int64, cfg *pack.ProactiveConfig) error {
	pendingMsgID, ok, err := s.infra.Cache.AI.GetProactivePending(ctx, groupID)
	if err != nil {
		return fmt.Errorf("get proactive pending key: %w", err)
	}
	if !ok {
		return nil
	}
	if pendingMsgID != msgID {
		return nil
	}
	// 群聊可以继续保持安静；冷却和每日上限用于避免 bot 反复救场。
	if allowed, err := proactiveAllowed(ctx, s, groupID, botID, cfg); err != nil || !allowed {
		return err
	}

	// 使用最近聊天历史，而不是只看触发消息，让 bot 能自然接住真实话题。
	history, _, _, err := s.infra.DB.Chat.ListGroupMessages(ctx, groupID, proactiveHistoryLimit, 0)
	if err != nil {
		return fmt.Errorf("list recent group messages: %w", err)
	}

	reply, err := s.agent.ProcessProactiveMessage(ctx, botID, history)
	if err != nil {
		return fmt.Errorf("agent proactive message: %w", err)
	}
	if reply == "" {
		return nil
	}
	if err := sendGroupBotMessage(s, ctx, botID, groupID, reply); err != nil {
		return fmt.Errorf("send proactive reply: %w", err)
	}
	if err := markProactiveSent(ctx, s, groupID, botID, cfg); err != nil {
		return fmt.Errorf("mark proactive sent: %w", err)
	}
	return nil
}

// proactiveAllowed 检查主动发言频控。
// cooldown 控制同一 bot 在同一群的最小间隔，daily count 控制每日最多主动发言次数。
func proactiveAllowed(ctx context.Context, s *AIService, groupID, botID uint, cfg *pack.ProactiveConfig) (bool, error) {
	coolingDown, err := s.infra.Cache.AI.IsProactiveCoolingDown(ctx, groupID, botID)
	if err != nil {
		return false, fmt.Errorf("check proactive cooldown: %w", err)
	}
	if coolingDown {
		return false, nil
	}

	count, err := s.infra.Cache.AI.ProactiveDailyCount(ctx, groupID, botID)
	if err != nil {
		return false, fmt.Errorf("check proactive daily count: %w", err)
	}
	return count < cfg.MaxPerDay, nil
}

// markProactiveSent 在主动消息成功发送后记录频控状态。
// daily key 使用固定 TTL，避免跨天计数残留长期占用 Redis。
func markProactiveSent(ctx context.Context, s *AIService, groupID, botID uint, cfg *pack.ProactiveConfig) error {
	return s.infra.Cache.AI.MarkProactiveSent(
		ctx,
		groupID,
		botID,
		time.Duration(cfg.CooldownMinutes)*time.Minute,
		48*time.Hour,
	)
}

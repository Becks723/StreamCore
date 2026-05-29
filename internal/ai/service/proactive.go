package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"StreamCore/internal/pkg/pack"
	kitexai "StreamCore/kitex_gen/ai"
	"StreamCore/pkg/util"
	"github.com/redis/go-redis/v9"
)

const (
	proactiveHistoryLimit = 20
	proactivePendingTTL   = time.Hour
)

// ProactiveCheck schedules a quiet-chat check for a newly sent human group message.
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

	pendingKey := proactivePendingKey(groupID)
	quietDelay := time.Duration(cfg.QuietMinutes) * time.Minute
	if err := s.rdb.Set(s.ctx, pendingKey, strconv.FormatInt(req.MsgId, 10), quietDelay+proactivePendingTTL).Err(); err != nil {
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

func selectProactiveBot(ctx context.Context, s *AIService, groupID, fromUID uint) (uint, *pack.ProactiveConfig, bool, error) {
	bots, err := s.db.ListGroupBots(ctx, groupID)
	if err != nil {
		return 0, nil, false, fmt.Errorf("ProactiveCheck: list group bots: %w", err)
	}
	for _, bot := range bots {
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

func runProactiveCheck(ctx context.Context, s *AIService, groupID, botID uint, msgID int64, cfg *pack.ProactiveConfig) error {
	pendingMsgID, err := s.rdb.Get(ctx, proactivePendingKey(groupID)).Int64()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get proactive pending key: %w", err)
	}
	if pendingMsgID != msgID {
		return nil
	}
	if allowed, err := proactiveAllowed(ctx, s, groupID, botID, cfg); err != nil || !allowed {
		return err
	}

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

func proactiveAllowed(ctx context.Context, s *AIService, groupID, botID uint, cfg *pack.ProactiveConfig) (bool, error) {
	exists, err := s.rdb.Exists(ctx, proactiveCooldownKey(groupID, botID)).Result()
	if err != nil {
		return false, fmt.Errorf("check proactive cooldown: %w", err)
	}
	if exists > 0 {
		return false, nil
	}

	count, err := s.rdb.Get(ctx, proactiveDailyKey(groupID, botID)).Int64()
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("check proactive daily count: %w", err)
	}
	return count < cfg.MaxPerDay, nil
}

func markProactiveSent(ctx context.Context, s *AIService, groupID, botID uint, cfg *pack.ProactiveConfig) error {
	if err := s.rdb.Set(ctx, proactiveCooldownKey(groupID, botID), "1", time.Duration(cfg.CooldownMinutes)*time.Minute).Err(); err != nil {
		return err
	}
	dailyKey := proactiveDailyKey(groupID, botID)
	if err := s.rdb.Incr(ctx, dailyKey).Err(); err != nil {
		return err
	}
	return s.rdb.Expire(ctx, dailyKey, 48*time.Hour).Err()
}

func proactivePendingKey(groupID uint) string {
	return fmt.Sprintf("ai:proactive:pending:%d", groupID)
}

func proactiveCooldownKey(groupID, botID uint) string {
	return fmt.Sprintf("ai:proactive:last:%d:%d", groupID, botID)
}

func proactiveDailyKey(groupID, botID uint) string {
	return fmt.Sprintf("ai:proactive:daily:%s:%d:%d", time.Now().Format("20060102"), groupID, botID)
}

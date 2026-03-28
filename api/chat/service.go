package chat

import (
	"fmt"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/hertz-contrib/websocket"
)

type WSService struct {
	hub *connectionHub
}

var (
	wsServiceOnce sync.Once
	wsServiceInst *WSService
)

func GlobalWSService() *WSService {
	wsServiceOnce.Do(func() {
		wsServiceInst = &WSService{hub: newConnectionHub()}
	})
	return wsServiceInst
}

func (s *WSService) RegisterClient(uid uint, conn *websocket.Conn) {
	s.hub.register(uid, conn)
}

func (s *WSService) UnregisterClient(uid uint, conn *websocket.Conn) {
	s.hub.unregister(uid, conn)
}

func (s *WSService) PushToUser(uid uint, msgType string, data any) error {
	buf, err := sonic.Marshal(&MsgWrapper{Type: msgType, Data: mustMarshalRaw(data)})
	if err != nil {
		return err
	}
	conn, ok := s.hub.get(uid)
	if !ok {
		return nil
	}
	conn.writeMu.Lock()
	defer conn.writeMu.Unlock()

	if err = conn.conn.WriteMessage(websocket.TextMessage, buf); err != nil {
		s.hub.unregister(uid, conn.conn)
		_ = conn.conn.Close()
		return fmt.Errorf("ws write failed uid=%d: %w", uid, err)
	}
	return nil
}

func mustMarshalRaw(v any) []byte {
	buf, err := sonic.Marshal(v)
	if err != nil {
		return []byte("{}")
	}
	return buf
}

package chat

import (
	"sync"

	"github.com/hertz-contrib/websocket"
)

type wsClient struct {
	conn    *websocket.Conn
	writeMu sync.Mutex
}

type connectionHub struct {
	mu      sync.RWMutex
	clients map[uint]*wsClient
}

func newConnectionHub() *connectionHub {
	return &connectionHub{
		clients: make(map[uint]*wsClient),
	}
}

func (h *connectionHub) register(uid uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if oldClient, ok := h.clients[uid]; ok && oldClient != nil && oldClient.conn != nil && oldClient.conn != conn {
		_ = oldClient.conn.Close()
	}
	h.clients[uid] = &wsClient{conn: conn}
}

func (h *connectionHub) unregister(uid uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if cur, ok := h.clients[uid]; ok && cur.conn == conn {
		delete(h.clients, uid)
	}
}

func (h *connectionHub) get(uid uint) (*wsClient, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.clients[uid]
	return client, ok
}

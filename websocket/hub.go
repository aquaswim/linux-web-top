package websocket

import (
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"sync"
)

type Hub struct {
	notEmpty bool
	conns    sync.Map
}

func NewHub() Hub {
	return Hub{
		conns:    sync.Map{},
		notEmpty: false,
	}
}

func (h *Hub) Broadcast(msg interface{}) int {
	size := 0
	h.conns.Range(func(key, value any) bool {
		//k := key.(uuid.UUID)
		v := value.(*websocket.Conn)
		size++
		err := v.WriteJSON(msg)
		if err != nil {
			return false
		}
		return true
	})
	h.notEmpty = size > 0
	return size
}

func (h *Hub) Register(conn *websocket.Conn) uuid.UUID {
	id := uuid.New()
	h.conns.Store(id, conn)
	h.notEmpty = true
	return id
}

func (h *Hub) Unregister(id uuid.UUID) {
	h.conns.Delete(id)
}

func (h *Hub) IsNotEmpty() bool {
	return h.notEmpty
}

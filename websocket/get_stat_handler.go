package websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
)

func NewGetStatHandler(hub *Hub) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		id := hub.Register(conn)
		log.Println("New client connected with id", id.String())
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				} else {
					log.Println("connection closed", err)
				}
				// return the loop since connection is closed
				log.Println("Unregister client", id.String())
				hub.Unregister(id)
				return
			}
		}
	})
}

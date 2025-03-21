package middleware

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type wsHandler func(conn *websocket.Conn, r *http.Request)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsMiddleware(next wsHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error upgrading connection: %v", err)
			return
		}
		defer wsConn.Close()

		next(wsConn, r)
	}
}

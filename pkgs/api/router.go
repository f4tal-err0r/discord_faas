package api

import (
	"net/http"

	"github.com/f4tal-err0r/discord_faas/pkgs/discord"
	"github.com/f4tal-err0r/discord_faas/pkgs/security"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type wsHandler func(conn *websocket.Conn, r *http.Request)

type Router struct {
	Router *mux.Router
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsMiddleware(next wsHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer wsConn.Close()

		next(wsConn, r)
	}
}

func NewRouter(bot *discord.Client, jwtSvc *security.JWTService) *Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/deploy", wsMiddleware(DeployHandler))
	r.HandleFunc("/api/context", ContextHandlerFunc(bot.ContextTokenCache, bot.Session))
	r.HandleFunc("/api/context/auth", AuthHandlerFunc(jwtSvc))

	return &Router{Router: r}
}

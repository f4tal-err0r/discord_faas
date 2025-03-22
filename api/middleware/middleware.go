package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/f4tal-err0r/discord_faas/pkgs/security"
	"github.com/gorilla/websocket"
)

var publicKey string

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

type JWTMiddleware struct {
	jwtsvc *security.JWTService
}

func NewJWTMiddleware(jwtsvc *security.JWTService) *JWTMiddleware {
	return &JWTMiddleware{jwtsvc: jwtsvc}
}

func (j *JWTMiddleware) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized: Token is missing", http.StatusUnauthorized)
			return
		}

		err := j.jwtsvc.VerifyToken(token)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		claims, err := j.jwtsvc.ParseToken(token)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

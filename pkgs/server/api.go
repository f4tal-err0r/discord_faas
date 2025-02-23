package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	pb "github.com/f4tal-err0r/discord_faas/proto"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

// Define an upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket handler type
type wsHandler func(conn *websocket.Conn, r *http.Request)

// Middleware to upgrade all connections to WebSocket
func wsWrapper(next wsHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the connection to a WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Connection Upgrade error:", err)
			return
		}
		defer conn.Close()

		// Call the next handler with the WebSocket connection
		next(conn, r)
	}
}

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/deploy", wsWrapper(DeployHandler))
	router.HandleFunc("/api/context", ContextHandler)
	router.HandleFunc("/api/context/auth", AuthHandler)
	router.HandleFunc("/api/context/decode", func(w http.ResponseWriter, r *http.Request) {
		token := getTokenFromHeader(r)
		if token == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, err := jwtsvc.ParseToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("%+v", claims)))

	})
	return router
}

func ContextHandler(w http.ResponseWriter, r *http.Request) {

	cfg, err := config.New()
	if err != nil {
		log.Println("Error getting config:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := r.Header.Get("Token")

	guildid, ok := contextToken.Load(token)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	} else {
		contextToken.Delete(r.URL.Query().Get("token"))
	}

	guild, err := GetGuildInfo(dsession, guildid.(string))
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			http.Error(w, "Guild not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "429") {
			http.Error(w, "Rate Limit for SelfGuildInfo exceeded", http.StatusTooManyRequests)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctxPb, err := proto.Marshal(&pb.ContextResp{
		ClientID:  cfg.Discord.ClientID,
		GuildID:   guild.ID,
		GuildName: guild.Name,
	})
	if err != nil {
		log.Println("Error marshalling ctx: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(ctxPb)
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromHeader(r)
	if token == "" {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	}

	guildid := r.Header.Get("GuildID")

	gm, err := FetchGuildMember(token, guildid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jwt, err := jwtsvc.CreateToken(Claims{
		UserID:  gm.User.ID,
		GuildID: guildid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)),
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(jwt))
}

func DeployHandler(conn *websocket.Conn, r *http.Request) {
	// Recieve a file from the request
	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		conn.Close()
		return
	}

	// Assume the body is a file
	_ = bytes.NewReader(p)
}

// func unmarshalRequest(r *http.Request) (*pb.Wrapper, error) {
// 	wrapper := new(pb.Wrapper)

// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid request: %w", err)
// 	}
// 	if err := proto.Unmarshal(data, wrapper); err != nil {
// 		return nil, fmt.Errorf("invalid request: %w", err)
// 	}

// 	return wrapper, nil
// }

func getTokenFromHeader(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		return ""
	}

	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return ""
	}

	return token
}

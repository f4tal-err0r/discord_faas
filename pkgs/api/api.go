package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	pb "github.com/f4tal-err0r/discord_faas/proto"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	"github.com/f4tal-err0r/discord_faas/pkgs/discord"
	"github.com/f4tal-err0r/discord_faas/pkgs/security"
)

func ContextHandler(w http.ResponseWriter, r *http.Request, contextTokenCache *sync.Map, dsession *discordgo.Session) {

	cfg, err := config.New()
	if err != nil {
		log.Println("Error getting config:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := r.Header.Get("Token")

	guildid, ok := contextTokenCache.Load(token)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	} else {
		contextTokenCache.Delete(r.URL.Query().Get("token"))
	}

	guild, err := dsession.Guild(guildid.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

func AuthHandler(w http.ResponseWriter, r *http.Request, jwtService *security.JWTService) {
	token := getTokenFromHeader(r)
	if token == "" {
		http.Error(w, "Unauthorized: Token is missing", http.StatusUnauthorized)
		return
	}

	guildID := r.Header.Get("GuildID")

	guildMember, err := discord.IdentGuildMember(token, guildID)
	if err != nil {
		http.Error(w, "Error identifying guild member: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if (guildMember.Permissions & (1 << 3)) != 0 {
		http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		return
	}

	jwtToken, err := jwtService.CreateToken(security.Claims{
		UserID:  guildMember.User.ID,
		GuildID: guildID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)),
		},
	})
	if err != nil {
		http.Error(w, "Error creating JWT: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(jwtToken))
}

func DeployHandler(conn *websocket.Conn, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		log.Println(err)
		return
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return
		}
		//send message to websocket
		err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Recieved file: %s", part.FileName())))
		if err != nil {
			log.Println(err)
			return
		}

	}

}

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

func ContextHandlerFunc(tokenCache *sync.Map, session *discordgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ContextHandler(w, r, tokenCache, session)
	}
}

func AuthHandlerFunc(jwtSvc *security.JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AuthHandler(w, r, jwtSvc)
	}
}

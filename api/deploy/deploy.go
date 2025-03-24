package deploy

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/protobuf/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/security"
	pb "github.com/f4tal-err0r/discord_faas/proto"
	"github.com/gorilla/mux"
)

func (h *Handler) DeployHandler(w http.ResponseWriter, r *http.Request) {
	var BuildReq pb.BuildFunc
	var ccs *discordgo.ApplicationCommand

	//get claims from context
	claims := r.Context().Value("claims").(*security.Claims)
	gid := claims.GuildID

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
		if part.FormName() == "metadata" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			err = proto.Unmarshal(buf.Bytes(), &BuildReq)
			if err != nil {
				log.Println(err)
				return
			}

			for _, c := range BuildReq.Commands {
				ccs, err = h.dbot.AddGuildCommand(c, gid)
				if err != nil {
					log.Println(err)
					w.Write([]byte("Error creating command: " + err.Error()))
				}
			}
		}
		if part.FileName() == "func" {
			cmdid := generateHex()
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			file, err := os.Create(cmdid + ".tar.gz")
			if err != nil {
				log.Println(err)
				return
			}
			_, err = io.Copy(buf, file)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
	// Fetch the guild data
	guild, err := h.dbot.Session.State.Guild(ccs.GuildID)
	if err != nil {
		// If the guild is not found in cache, fetch it from Discord API
		guild, err = h.dbot.Session.Guild(ccs.GuildID)
		if err != nil {
			log.Fatal("error retrieving guild:", err)
			return
		}
	}

	w.Write([]byte(fmt.Sprintf("%s Function deployed successfully to %s", BuildReq.GetName(), guild.Name)))
}

func (h *Handler) AddRoute(r *mux.Router) {
	r.HandleFunc("/api/func/deploy", h.DeployHandler)
}

func (h *Handler) IsSecure() bool {
	return true
}

func generateHex() string {
	b := make([]byte, 5) // 5 bytes = 10 hex characters
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

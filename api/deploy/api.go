package deploy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"google.golang.org/protobuf/proto"

	pb "github.com/f4tal-err0r/discord_faas/proto"
	"github.com/gorilla/mux"
)

func (h *Handler) DeployHandler(w http.ResponseWriter, r *http.Request) {
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
			var BuildReq pb.BuildFunc
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			err = proto.Unmarshal(buf.Bytes(), &BuildReq)
			if err != nil {
				log.Println(err)
				return
			}

			w.Write([]byte(fmt.Sprintf("Recieved metadata: %+v", BuildReq)))
			break
		}

		_, err = w.Write([]byte(fmt.Sprintf("Recieved file: %s", part.FileName())))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (h *Handler) AddRoute(r *mux.Router) {
	r.HandleFunc("/api/deploy", h.DeployHandler)
}

func (h *Handler) IsSecure() bool {
	return true
}

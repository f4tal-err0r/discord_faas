package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	pb "github.com/f4tal-err0r/discord_faas/proto"
	fzf "github.com/ktr0731/go-fuzzyfinder"
	"google.golang.org/protobuf/proto"
)

//TODO: Serialize future JWT token to server here, verifying ident w/ Oauth token

func NewContext(url string, guildid string) *pb.ContextResp {
	req, err := http.NewRequest("POST", url+"/api/context", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Send the request with headers
	oauth, err := NewUserAuth().GetToken()
	if err != nil {
		log.Fatal(err)
	}

	contextRequest := pb.GetContext{
		GuildID: guildid,
		Token:   oauth,
	}

	// Send protobuf request to server

	msgBody, err := proto.Marshal(&contextRequest)
	if err != nil {
		log.Fatal(err)
	}

	//write msgBody to req.Body as protobuf bytes
	req.Body = io.NopCloser(bytes.NewReader(msgBody))
	req.Header.Set("Content-Type", "application/x-protobuf")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var ctx pb.ContextResp
	// Decode the response
	if err := json.NewDecoder(resp.Body).Decode(&ctx); err != nil {
		fmt.Print("Unable to decode response")
		log.Fatal(err)
	}

	CtxList, err := LoadContextList()
	if err != nil {
		fmt.Print("Unable to load context list")
		log.Fatal(err)
	}
	for _, ctx := range CtxList {
		ctx.CurrentContext = false
	}

	ctx.CurrentContext = true

	//Append ctx to ContextList only if guildid is not already present
	for _, ctxl := range CtxList {
		if ctxl.GuildID == guildid {
			fmt.Printf("%s selected\n", ctx.GuildName)
			return &ctx
		}
	}

	CtxList = append(CtxList, &ctx)
	err = SerializeContextList(CtxList)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s Added Successfully and selected\n", ctx.GuildName)

	return &ctx
}

func SerializeContextList(ctxl []*pb.ContextResp) error {
	cacheDir := FetchCacheDir("context")

	file, err := createFileIfNotExists(cacheDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// If error returns EOF, return empty list
	if err == io.EOF {
		return nil
	}

	return json.NewEncoder(file).Encode(ctxl)
}

// Load context from cache
func LoadContextList() ([]*pb.ContextResp, error) {
	var localctx []*pb.ContextResp
	cacheDir := FetchCacheDir("context")
	file, err := createFileIfNotExists(cacheDir)
	// If error returns EOF, return empty list
	if err != nil && err == io.EOF {
		return localctx, nil
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&localctx)
	if err != nil {
		// If error returns EOF, return empty list
		if err == io.EOF {
			return localctx, nil
		}
		return nil, err
	}
	return localctx, nil
}

func SwitchContext(ctxl []*pb.ContextResp, gid string) {
	if len(ctxl) == 0 {
		fmt.Println("No contexts found")
	}
	for _, ctx := range ctxl {
		if ctx.CurrentContext {
			ctx.CurrentContext = false
		} else if ctx.GuildID == gid {
			ctx.CurrentContext = true
		}
	}
}

func GetCurrentContext() *pb.ContextResp {
	ctxl, err := LoadContextList()
	if err != nil {
		log.Fatal(err)
	}
	for _, ctx := range ctxl {
		if ctx.CurrentContext {
			return ctx
		}
	}
	return nil
}

func ListContexts() {
	ContextList, err := LoadContextList()
	if err != nil {
		log.Fatalf("failed to load context list: %v", err)
	}

	if len(ContextList) == 0 {
		fmt.Println("No context found")
		return
	}

	_, err = fzf.Find(ContextList, func(i int) string {
		SwitchContext(ContextList, ContextList[i].GuildID)
		return ContextList[i].GuildName
	},
	)
	if err != nil {
		log.Fatalf("failed to select context: %v", err)
	}
}

func createFileIfNotExists(filePath string) (*os.File, error) {
	// os.O_CREATE creates the file if it does not exist
	// os.O_EXCL returns an error if the file already exists
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		if os.IsExist(err) {
			// If the file already exists, open it in read-write mode
			file, err = os.OpenFile(filePath, os.O_RDWR, 0666)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return file, nil
}

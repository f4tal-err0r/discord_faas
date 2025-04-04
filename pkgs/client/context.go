package client

import (
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

func NewContext(uri string, token string) *pb.ContextResp {
	req, err := http.NewRequest(http.MethodGet, uri+"/api/context", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var ctx pb.ContextResp

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Unable to get context: ", resp.Status, string(body))
	}

	// Decode the response
	if err := proto.Unmarshal(body, &ctx); err != nil {
		fmt.Print("Unable to decode response")
		log.Fatal(err)
	}

	CtxList, err := LoadContextList()
	if err != nil {
		fmt.Print("Unable to load context list")
		log.Fatal(err)
	}

	ctx.CurrentContext = true
	ctx.ServerURL = uri

	//Append ctx to ContextList only if guildid is not already present
	for _, c := range CtxList {
		if c.GuildID == ctx.GuildID {
			fmt.Printf("%s already exists, selected as current context\n", ctx.GuildName)
			return &ctx
		}
	}
	for _, ctxStore := range CtxList {
		ctxStore.CurrentContext = false
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

func GetCurrentContext() (*pb.ContextResp, error) {
	ctxl, err := LoadContextList()
	if err != nil {
		log.Fatal(err)
	}
	for _, ctx := range ctxl {
		if ctx.CurrentContext {
			return ctx, nil
		}
	}
	return nil, fmt.Errorf("no current context found")
}

func UpdateContextToken(ctx *pb.ContextResp) error {
	ctxl, err := LoadContextList()
	if err != nil {
		return err
	}
	for _, c := range ctxl {
		if c.GuildID == ctx.GuildID {
			c.JWToken = ctx.JWToken
		}
	}
	return SerializeContextList(ctxl)
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

func AuthContent(ctx *pb.ContextResp) error {
	oauth := NewUserAuth()

	oauthToken, err := oauth.StartAuth()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, ctx.ServerURL+"/api/context/auth", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", oauthToken.AccessToken))
	req.Header.Add("GuildID", ctx.GuildID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("err while authenticating to context: ", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Unable to authenticate to context: ", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Unable to read response body: ", err)
	}

	ctx.JWToken = string(body)

	if err := UpdateContextToken(ctx); err != nil {
		log.Fatal("Unable to update context token: ", err)
	}

	return nil

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

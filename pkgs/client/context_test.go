package client_test

import (
	"io"
	"os"
	"testing"

	"github.com/f4tal-err0r/discord_faas/pkgs/client"
)

// This only works locally with a valid oauth token
// func TestContext(t *testing.T) {
// 	// create a rest server
// 	router := mux.NewRouter()
// 	server := httptest.NewServer(router)

// 	// respond on /api/context
// 	router.HandleFunc("/api/context", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte(`{"client_id":"123456789012345678","guild_id":"123456789012345678","user_id":"123456789012345678"}`))
// 	})

// 	// Create a context
// 	ctx := client.NewContext(server.URL, "123456789012345678")

// 	// Check the context
// 	if ctx.ClientID != "123456789012345678" {
// 		t.Errorf("Unexpected client ID: %s", ctx.ClientID)
// 	}
// 	if ctx.GuildID != "123456789012345678" {
// 		t.Errorf("Unexpected guild ID: %s", ctx.GuildID)
// 	}
// }

func TestSerializeContextList(t *testing.T) {
	// Test case: Empty context list
	ctxList := []*client.ContextResp{}
	err := client.SerializeContextList(ctxList)
	if err != nil {
		t.Errorf("Expected no error for empty context list, got %v", err)
	}

	// Test case: Non-empty context list
	ctxList = []*client.ContextResp{
		{
			ClientID:  "client1",
			GuildID:   "guild1",
			GuildName: "Guild 1",
		},
		{
			ClientID:  "client2",
			GuildID:   "guild2",
			GuildName: "Guild 2",
		},
	}
	expected := `[{"ClientID":"client1","GuildID":"guild1","GuildName":"Guild 1"},{"ClientID":"client2","GuildID":"guild2","GuildName":"Guild 2"}]`
	err = client.SerializeContextList(ctxList)
	if err != nil {
		t.Errorf("Expected no error for non-empty context list, got %v", err)
	}
	file, err := os.Open(client.FetchCacheDir("context"))
	if err != nil {
		t.Errorf("Failed to open cache file: %v", err)
	}
	defer file.Close()
	actual, err := io.ReadAll(file)
	if err != nil {
		t.Errorf("Failed to read cache file: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("Expected serialized context list %s, got %s", expected, string(actual))
	}
}

func TestLoadContextList(t *testing.T) {
	// Test case: Empty context list
	ctxList, err := client.LoadContextList()
	if err != nil {
		t.Errorf("Expected no error for empty context list, got %v", err)
	}
	if len(ctxList) != 0 {
		t.Errorf("Expected empty context list, got %v", ctxList)
	}
}

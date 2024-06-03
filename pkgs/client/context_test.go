package client_test

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/f4tal-err0r/discord_faas/pkgs/client"
// 	"github.com/gorilla/mux"
// )

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

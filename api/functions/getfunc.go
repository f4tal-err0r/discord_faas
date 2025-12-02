package functions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/f4tal-err0r/discord_faas/pkgs/security"
)

func (h *Handler) GetFuncsHandler(w http.ResponseWriter, r *http.Request) {
	// get claims from context
	claims := r.Context().Value("claims").(security.Claims)
	if claims.GuildID == "" {
		http.Error(w, "GuildID missing in token claims", http.StatusBadRequest)
		return
	}

	//Marshall jwt claims to JSON and write to response
	resp, err := json.Marshal(claims)
	if err != nil {
		http.Error(w, "Error marshalling claims: "+err.Error(), http.StatusInternalServerError)
		fmt.Print("Error marshalling claims: " + err.Error())
		return
	}

	w.Write(resp)
}

package api

import (
	middleware "github.com/f4tal-err0r/discord_faas/api/middleware"
	"github.com/f4tal-err0r/discord_faas/pkgs/security"
	"github.com/gorilla/mux"
)

type RouterAdder interface {
	AddRoute(r *mux.Router)
	IsSecure() bool
}

func NewRouter(jwtsvc *security.JWTService, route ...RouterAdder) (*mux.Router, error) {
	r := mux.NewRouter()

	jwtmw := middleware.NewJWTMiddleware(jwtsvc)
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(jwtmw.JWTMiddleware)

	for _, v := range route {
		if v.IsSecure() {
			v.AddRoute(protected)
		} else {
			v.AddRoute(r)
		}
	}

	return r, nil
}

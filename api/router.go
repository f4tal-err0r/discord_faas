package api

import (
	"github.com/gorilla/mux"
)

type RouterAdder interface {
	AddRoute(r *mux.Router)
	IsSecure() bool
}

func NewRouter(route ...RouterAdder) (*mux.Router, error) {
	r := mux.NewRouter()

	for _, v := range route {
		if v.IsSecure() {
			//TODO: JWT Check
		}
		v.AddRoute(r)
	}

	return r, nil
}

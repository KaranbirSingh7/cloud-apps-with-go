package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type pinger interface {
	Ping(ctx context.Context) error
}

func HealthzWrapper(r *mux.Router, p pinger) {
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := p.Ping(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		fmt.Fprint(w, "up")
	})
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "up")
}

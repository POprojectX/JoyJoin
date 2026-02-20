package httpdelivery

import (
	"log"
	"net/http"
	"time"
)

func Run(h *Handler, addr string) {
	srv := &http.Server{
		Addr: addr,
		Handler: h.Routes(),
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	log.Printf("Server running on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
package httpdelivery

import (
	"net/http"
	"user-service/internal/http/middleware"
)

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	authMiddleware := middleware.ParseJWT("DanilTopRonaldoTop")

	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	mux.Handle("/register", h.Register())
	mux.Handle("/login", h.Login())
	mux.Handle("/googleLogin", h.LoginWithGoogle())
	
	mux.Handle("/me", authMiddleware(h.GetMe()))
	return mux
}
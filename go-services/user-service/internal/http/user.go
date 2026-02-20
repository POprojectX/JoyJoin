package httpdelivery

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, jwt, err := h.Orchestrator.RegisterAndCreateProfile(r.Context(), req.Email, req.Password, req.FirstName, req.LastName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responce := map[string]interface{}{
			"user":  user,
			"token": jwt,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(responce)
	}
}
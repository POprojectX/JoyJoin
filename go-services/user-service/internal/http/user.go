package httpdelivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return 
		}
		
		dto, jwt, err := h.Orchestrator.LoginAndGetProfile(r.Context(), req.Email, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
            return
		}
		response := map[string]interface{}{
            "user":  dto,
            "token": jwt,
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
	}
}

func (h *Handler) LoginWithGoogle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Token string `json:"token"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
            return
		}

		dto, jwt, err := h.Orchestrator.LoginWithGoogleAndGetProfile(r.Context(), req.Token)
		if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

		json.NewEncoder(w).Encode(map[string]interface{} {
			"user": dto,
			"token": jwt,
		})
	}
}

func(h *Handler) GetMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value("userID")
		userIDStr, ok := val.(string)
		if !ok {
			http.Error(w, "user not found in context", http.StatusUnauthorized)
			return 
		}
		
		uuidUser, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "invalid uuid format", http.StatusBadRequest)
			return 
		}

		user, err := h.Orchestrator.GetUserFullProfile(r.Context(), uuidUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}

		responce := map[string]interface{} {
			"id": user.ID,
			"email": user.Email,
			"first_name": user.FirstName,
			"last_name": user.LastName, 
			"events_count": user.EventsCount,
			"participants": user.Participations,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
	}
}

func (h *Handler) GetFullUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID string `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return 
		}
		uuidUser, err := uuid.Parse(req.UserID)
		if err != nil {
			http.Error(w, "invalid uuid format", http.StatusBadRequest)
		}
		user, err := h.Orchestrator.GetUserFullProfile(r.Context(), uuidUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}

		responce := map[string]interface{} {
			"id": user.ID,
			"email": user.Email,
			"first_name": user.FirstName,
			"last_name": user.LastName, 
			"events_count": user.EventsCount,
			"participants": user.Participations,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
	}
}

func (h *Handler) DeleteAndCleanUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID string `json:"user_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return 
		}

		uuidUser, err := uuid.Parse(req.UserID)
		if err != nil {
			http.Error(w, "invalid uuid format", http.StatusBadRequest)
		}

		user := h.Orchestrator.DeleteUserAndCleanup(r.Context(), uuidUser)
		if user != nil {
			http.Error(w, user.Error(), http.StatusInternalServerError)
			return
		}

		responce := fmt.Sprintf("user: %s was deleted successfully!",req.UserID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
	}
}
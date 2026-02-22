package httpdelivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) GetEventFullDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			EventID uuid.UUID `json:"event_id"`
			RequesterID uuid.UUID `json:"requester_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		role, participants, err := h.Orchestrator.GetEventFullDetails(r.Context(), req.EventID, req.RequesterID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		responce := map[string]interface{} {
			"role" : role,
			"participants" : participants,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
	}
}

func (h *Handler) TransferOwnership() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			CurrentID uuid.UUID `json:"current_id"`
			NewOwnerID uuid.UUID `json:"new_owner_id"`
			EventID uuid.UUID `json:"event_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		transfer := h.Orchestrator.TransferOwnership(r.Context(), req.CurrentID, req.NewOwnerID, req.EventID)
		if transfer != nil {
			http.Error(w, transfer.Error(), http.StatusInternalServerError)
			return 
		}

		responce := fmt.Sprintf("Transfer ownership for event: %v finish successfully! New Owner: %v", req.EventID, req.NewOwnerID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
	}
}
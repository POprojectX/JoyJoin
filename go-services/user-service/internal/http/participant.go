package httpdelivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) JoinEventAsGuestPublic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID uuid.UUID `json:"user_id"`
			EventID uuid.UUID `json:"event_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		join := h.Orchestrator.JoinEventAsGuestPublic(r.Context(), req.UserID, req.EventID)
		if join != nil {
			http.Error(w, join.Error(), http.StatusInternalServerError)
		}

		responce := fmt.Sprintf("user successfully join to event: %v as a Guest", req.EventID)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
	}
}

func (h *Handler) JoinEventAsGuestPrivate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			RequesterID uuid.UUID `json:"requester_id"`
			UserID uuid.UUID `json:"user_id"`
			EventID uuid.UUID `json:"event_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		join := h.Orchestrator.JoinEventAsGuestPrivate(r.Context(), req.RequesterID, req.UserID, req.EventID)
		if join != nil {
			http.Error(w, join.Error(), http.StatusInternalServerError)
		}

		responce := fmt.Sprintf("user successfully join to private event: %v as a Guest", req.EventID)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
	}
}

func (h *Handler) LeaveEventAndFreeSlot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID uuid.UUID `json:"user_id"`
			EventID uuid.UUID `json:"event_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		join := h.Orchestrator.LeaveEventAndFreeSlot(r.Context(), req.UserID, req.EventID)
		if join != nil {
			http.Error(w, join.Error(), http.StatusInternalServerError)
		}

		responce := fmt.Sprintf("user successfully leave event: %v", req.EventID)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
	}
}

func (h *Handler) AssignStaff() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			RequesterID uuid.UUID `json:"requester_id"`
			TargetUserID uuid.UUID `json:"target_user_id"`
			EventID uuid.UUID `json:"event_id"`
			ProfRole uint `json:"prof_role"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		assignStaff := h.Orchestrator.AssignStaffWithOutSlotCheck(r.Context(), req.RequesterID, req.TargetUserID, req.EventID, &req.ProfRole)
		if assignStaff != nil {
			http.Error(w, assignStaff.Error(), http.StatusInternalServerError)
		}

		responce := fmt.Sprintf("user: %v successfully assign to event: %v with IDrole: %v",req.TargetUserID, req.EventID, req.ProfRole)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
	}
}
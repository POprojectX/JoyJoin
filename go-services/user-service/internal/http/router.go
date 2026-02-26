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

	// user routes
	mux.Handle("/register", h.Register())
	mux.Handle("/login", h.Login())
	mux.Handle("/googleLogin", h.LoginWithGoogle())
	mux.Handle("/Me", h.GetMe())
	mux.Handle("/getfulluser", h.GetFullUser())
	mux.Handle("/deleteUser", h.DeleteAndCleanUp())
	
	//event routes
	mux.Handle("/createEvent", h.CreateEventWithOwner())
	mux.Handle("/publishEvent", h.PublishEventAtomic())
	mux.Handle("/draftEvent", h.DraftEvent())
	mux.Handle("/cancelEvent", h.CancelEventWithCleanup())
	mux.Handle("/deleteEvent", h.DeleteEventWithPermissions())
	mux.Handle("/getEventsBuUserSlot", h.GetEventByUserRole())
	mux.Handle("/getEventsByLocation", h.GetEventsByLocatoin())
	mux.Handle("/hasSlots", h.HasAailableSlots())

	//participant routes
	mux.Handle("/joinEventAsGuestPublic", h.JoinEventAsGuestPublic())
	mux.Handle("/joinEventAsGuestPrivate", h.JoinEventAsGuestPrivate())
	mux.Handle("/leaveEventAsGuest", h.LeaveEventAndFreeSlot())
	mux.Handle("/assignStaff", h.AssignStaff())
	//mux.Handle("/")

	mux.Handle("/me", authMiddleware(h.GetMe()))
	return mux
}
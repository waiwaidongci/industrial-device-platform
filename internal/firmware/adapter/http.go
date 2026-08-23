package adapter

import (
	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/firmware/application"
	"github.com/example/industrial-device-platform/internal/firmware/domain"
	"net/http"
)

type Handler struct{ svc *application.Service }

func NewHandler(s *application.Service) *Handler { return &Handler{svc: s} }
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var x domain.Release
		if e := internal.Decode(r, &x); e != nil {
			internal.JSON(w, 400, map[string]string{"error": e.Error()})
			return
		}
		if e := h.svc.Publish(r.Context(), &x); e != nil {
			internal.JSON(w, 400, map[string]string{"error": e.Error()})
			return
		}
		internal.JSON(w, 201, x)
		return
	}
	http.NotFound(w, r)
}

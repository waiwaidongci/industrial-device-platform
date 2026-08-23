package adapter

import (
	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/tenant/application"
	"github.com/example/industrial-device-platform/internal/tenant/domain"
	"net/http"
)

type Handler struct{ svc *application.Service }

func NewHandler(s *application.Service) *Handler { return &Handler{svc: s} }
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var t domain.Tenant
	if e := internal.Decode(r, &t); e != nil {
		internal.JSON(w, 400, map[string]string{"error": e.Error()})
		return
	}
	if e := h.svc.Create(r.Context(), &t); e != nil {
		internal.JSON(w, 400, map[string]string{"error": e.Error()})
		return
	}
	internal.JSON(w, 201, t)
}

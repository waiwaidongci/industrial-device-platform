package adapter

import (
	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/remotecommand/application"
	"github.com/example/industrial-device-platform/internal/remotecommand/domain"
	"net/http"
	"strings"
)

type Handler struct {
	svc     *application.Service
	metrics *internal.Metrics
}

func NewHandler(s *application.Service, m *internal.Metrics) *Handler {
	return &Handler{svc: s, metrics: m}
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(p) >= 2 && p[1] == "commands" {
		if r.Method == http.MethodPost && len(p) == 2 {
			var c domain.Command
			if e := internal.Decode(r, &c); e != nil {
				internal.JSON(w, 400, map[string]string{"error": e.Error()})
				return
			}
			c.IdempotencyKey = r.Header.Get("Idempotency-Key")
			out, e := h.svc.Submit(r.Context(), &c)
			if e != nil {
				internal.JSON(w, 409, map[string]string{"error": e.Error()})
				return
			}
			h.metrics.IncCommand()
			internal.JSON(w, 202, out)
			return
		}
		if r.Method == http.MethodGet && len(p) == 3 {
			out, e := h.svc.Get(r.Context(), p[2])
			if e != nil {
				internal.JSON(w, 404, map[string]string{"error": e.Error()})
				return
			}
			internal.JSON(w, 200, out)
			return
		}
		if r.Method == http.MethodGet && len(p) == 2 {
			out, _ := h.svc.List(r.Context(), r.URL.Query().Get("device_id"))
			internal.JSON(w, 200, out)
			return
		}
	}
	http.NotFound(w, r)
}

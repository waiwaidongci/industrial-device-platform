package adapter

import (
	"fmt"
	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/devicecatalog/application"
	"github.com/example/industrial-device-platform/internal/devicecatalog/domain"
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
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) >= 2 && parts[1] == "devices" {
		if r.Method == http.MethodPost && len(parts) == 2 {
			h.create(w, r)
			return
		}
		if len(parts) >= 3 {
			h.item(w, r, parts[2])
			return
		}
		if r.Method == http.MethodGet {
			h.list(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var d domain.Device
	if e := internal.Decode(r, &d); e != nil {
		internal.JSON(w, 400, map[string]string{"error": e.Error()})
		return
	}
	if e := h.svc.Register(r.Context(), &d); e != nil {
		internal.JSON(w, 409, map[string]string{"error": e.Error()})
		return
	}
	h.metrics.SetDevices(lenMust(h.svc.List(r.Context(), "", "")))
	internal.JSON(w, 201, d)
}
func lenMust(v []*domain.Device, e error) int {
	if e != nil {
		return 0
	}
	return len(v)
}
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	v, e := h.svc.List(r.Context(), r.URL.Query().Get("tenant_id"), r.URL.Query().Get("group"))
	if e != nil {
		internal.JSON(w, 500, map[string]string{"error": e.Error()})
		return
	}
	internal.JSON(w, 200, v)
}
func (h *Handler) item(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method == http.MethodGet {
		d, e := h.svc.Get(r.Context(), id)
		if e != nil {
			internal.JSON(w, 404, map[string]string{"error": e.Error()})
			return
		}
		internal.JSON(w, 200, d)
		return
	}
	if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/heartbeat") {
		d, e := h.svc.Heartbeat(r.Context(), id)
		if e != nil {
			internal.JSON(w, 404, map[string]string{"error": e.Error()})
			return
		}
		internal.JSON(w, 200, d)
		return
	}
	if r.Method == http.MethodPatch {
		var p struct {
			Status     domain.Status     `json:"status"`
			Name       string            `json:"name"`
			Group      string            `json:"group"`
			Tags       map[string]string `json:"tags"`
			Parameters map[string]any    `json:"parameters"`
		}
		if e := internal.Decode(r, &p); e != nil {
			internal.JSON(w, 400, map[string]string{"error": e.Error()})
			return
		}
		d, e := h.svc.Get(r.Context(), id)
		if e != nil {
			internal.JSON(w, 404, map[string]string{"error": e.Error()})
			return
		}
		if p.Name != "" {
			d.Name = p.Name
		}
		if p.Group != "" {
			d.Group = p.Group
		}
		if p.Status != "" {
			d.Status = p.Status
		}
		if p.Tags != nil {
			d.Tags = p.Tags
		}
		if p.Parameters != nil {
			d.Parameters = p.Parameters
		}
		e = h.svc.Update(r.Context(), d)
		if e != nil {
			internal.JSON(w, 400, map[string]string{"error": fmt.Sprint(e)})
			return
		}
		internal.JSON(w, 200, d)
		return
	}
	http.NotFound(w, r)
}

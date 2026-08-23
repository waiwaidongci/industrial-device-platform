package adapter

import (
	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/alert/application"
	"github.com/example/industrial-device-platform/internal/alert/domain"
	"net/http"
	"strings"
)

type Handler struct{ eval *application.Evaluator }

func NewHandler(e *application.Evaluator) *Handler { return &Handler{eval: e} }
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if r.Method == http.MethodPost && len(p) == 3 && p[1] == "alert-rules" {
		var x domain.Rule
		if e := internal.Decode(r, &x); e != nil {
			internal.JSON(w, 400, map[string]string{"error": e.Error()})
			return
		}
		if e := h.eval.Put(x); e != nil {
			internal.JSON(w, 400, map[string]string{"error": e.Error()})
			return
		}
		internal.JSON(w, 201, x)
		return
	}
	http.NotFound(w, r)
}

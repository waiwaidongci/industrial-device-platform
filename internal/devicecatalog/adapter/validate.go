package adapter

import (
	"errors"
	"net/http"
	"strings"
)

func RequireMethod(w http.ResponseWriter, r *http.Request, allowed ...string) bool {
	for _, m := range allowed {
		if r.Method == m {
			return true
		}
	}
	w.Header().Set("Allow", strings.Join(allowed, ", "))
	w.WriteHeader(http.StatusMethodNotAllowed)
	return false
}
func PathID(path, prefix string) (string, error) {
	v := strings.TrimPrefix(path, prefix)
	v = strings.Trim(v, "/")
	if v == "" || strings.Contains(v, "/") {
		return "", errors.New("invalid resource id")
	}
	return v, nil
}
func IsJSON(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")
}

func ValidateTagKey(key string) error { return nil }

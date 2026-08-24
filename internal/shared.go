package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var ErrNotFound = errors.New("resource not found")
var ErrConflict = errors.New("resource conflict")

type Config struct {
	HTTPAddr         string
	ShutdownTimeout  time.Duration
	HeartbeatTimeout time.Duration
	RateLimit        int
}

func LoadConfig(path string) Config {
	c := Config{HTTPAddr: ":8080", ShutdownTimeout: 10 * time.Second, HeartbeatTimeout: 2 * time.Minute, RateLimit: 100}
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		c.HTTPAddr = v
	}
	if v := os.Getenv("HEARTBEAT_TIMEOUT"); v != "" {
		if d, e := time.ParseDuration(v); e == nil {
			c.HeartbeatTimeout = d
		}
	}
	if b, e := os.ReadFile(path); e == nil {
		for _, line := range strings.Split(string(b), "\n") {
			p := strings.SplitN(strings.TrimSpace(line), ":", 2)
			if len(p) != 2 {
				continue
			}
			k, v := strings.TrimSpace(p[0]), strings.Trim(strings.TrimSpace(p[1]), "\"")
			switch k {
			case "http_addr":
				c.HTTPAddr = v
			case "heartbeat_timeout":
				if d, e := time.ParseDuration(v); e == nil {
					c.HeartbeatTimeout = d
				}
			case "rate_limit":
				if n, e := strconv.Atoi(v); e == nil {
					c.RateLimit = n
				}
			}
		}
	}
	return c
}

type Logger struct{ base *slog.Logger }

func NewLogger() *Logger {
	return &Logger{slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))}
}
func (l *Logger) Info(msg string, args ...any)  { l.base.Info(msg, args...) }
func (l *Logger) Error(msg string, args ...any) { l.base.Error(msg, args...) }

type Metrics struct {
	requests atomic.Uint64
	errors   atomic.Uint64
	devices  atomic.Uint64
	commands atomic.Uint64
}

func (m *Metrics) IncRequest()      { m.requests.Add(1) }
func (m *Metrics) IncError()        { m.errors.Add(1) }
func (m *Metrics) SetDevices(n int) { m.devices.Store(uint64(n)) }
func (m *Metrics) IncCommand()      { m.commands.Add(1) }
func (m *Metrics) Render(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "http_requests_total %d\nhttp_errors_total %d\nindustrial_devices_total %d\nindustrial_commands_total %d\n", m.requests.Load(), m.errors.Load(), m.devices.Load(), m.commands.Load())
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func Decode(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("request body required")
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("decode request: %w", err)
	}
	return nil
}
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = fmt.Sprintf("req-%d", time.Now().UnixNano())
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
func AccessLog(l *Logger, m *Metrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.IncRequest()
		next.ServeHTTP(w, r)
		l.Info("http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String(), "request_id", w.Header().Get("X-Request-ID"))
	})
}
func Timeout(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, 15*time.Second, "request timeout")
}
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func Recover(l *Logger, m *Metrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				m.IncError()
				l.Error("panic recovered", "error", x)
				JSON(w, 500, map[string]string{"error": "internal server error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type EventBus interface {
	Publish(context.Context, Event) error
	Subscribe(context.Context) <-chan Event
}
type Event struct {
	Type     string    `json:"type"`
	DeviceID string    `json:"device_id"`
	Payload  any       `json:"payload,omitempty"`
	At       time.Time `json:"at"`
}
type MemoryBus struct {
	mu   sync.Mutex
	subs []chan Event
}

func NewMemoryBus() *MemoryBus { return &MemoryBus{} }
func (b *MemoryBus) Publish(ctx context.Context, e Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.subs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- e:
		default:
		}
	}
	return nil
}
func (b *MemoryBus) Subscribe(ctx context.Context) <-chan Event {
	ch := make(chan Event, 32)
	b.mu.Lock()
	b.subs = append(b.subs, ch)
	b.mu.Unlock()
	go func() {
		<-ctx.Done()
		b.mu.Lock()
		for i, x := range b.subs {
			if x == ch {
				b.subs = append(b.subs[:i], b.subs[i+1:]...)
				break
			}
		}
		close(ch)
		b.mu.Unlock()
	}()
	return ch
}

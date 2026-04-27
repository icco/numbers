package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/icco/gutil/logging"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.uber.org/zap"
)

// serverName is the otelhttp span/metric scope.
const serverName = "numbers"

// Options configures the HTTP router. MetricsHandler is mounted at /metrics.
type Options struct {
	Logger         *zap.SugaredLogger
	MetricsHandler http.Handler
}

// Router returns the HTTP handler, wrapped with otelhttp (excluding /metrics).
func Router(opts Options) http.Handler {
	r := chi.NewRouter()
	r.Use(logging.Middleware(opts.Logger.Desugar()))
	r.Use(routeTag)

	r.Get("/", handleCharacter)
	r.Get("/json", handleJSON)
	r.Get("/healthz", handleHealthz)

	if opts.MetricsHandler != nil {
		r.Method(http.MethodGet, "/metrics", opts.MetricsHandler)
	}

	return otelhttp.NewHandler(r, serverName,
		otelhttp.WithFilter(func(req *http.Request) bool {
			return req.URL.Path != "/metrics"
		}),
	)
}

// routeTag stamps the chi route pattern onto otelhttp metric labels.
func routeTag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		labeler, ok := otelhttp.LabelerFromContext(r.Context())
		if !ok {
			return
		}
		if pattern := chi.RouteContext(r.Context()).RoutePattern(); pattern != "" {
			labeler.Add(semconv.HTTPRoute(pattern))
		}
	})
}

// handleCharacter writes the current character as a decimal codepoint.
func handleCharacter(w http.ResponseWriter, r *http.Request) {
	l := logging.FromContext(r.Context())
	d, err := GetCharacter(l)
	if err != nil {
		l.Errorw("get character", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := fmt.Fprintf(w, "%d", d.Character); err != nil {
		l.Debugw("write character", zap.Error(err))
	}
}

// handleJSON writes the current character lookup as JSON.
func handleJSON(w http.ResponseWriter, r *http.Request) {
	l := logging.FromContext(r.Context())
	d, err := GetCharacter(l)
	if err != nil {
		l.Errorw("get character", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(d); err != nil {
		l.Errorw("encode json", zap.Error(err))
	}
}

// handleHealthz is the liveness probe.
func handleHealthz(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("hi.")); err != nil {
		logging.FromContext(r.Context()).Errorw("write healthz", zap.Error(err))
	}
}

// Data holds a single character lookup result derived from the current time.
type Data struct {
	Character     rune  `json:"character"`
	SecondsPassed int64 `json:"seconds_passed"`
	Seconds       int64 `json:"seconds"`
	Length        int64 `json:"length"`
	Lookup        int64 `json:"lookup"`
}

// Log formats Data for human-readable debug output.
func (d *Data) Log() string {
	return fmt.Sprintf("(%v / %v) * %d = %d: %s (%d)", d.SecondsPassed, d.Seconds, d.Length, d.Lookup, string(d.Character), d.Character)
}

// GetCharacter returns the byte from book.txt that maps to the current
// position within the UTC week. The Lookup index is clamped because
// floating-point rounding can otherwise produce text[len(text)] in the
// final nanoseconds of the week.
func GetCharacter(l *zap.SugaredLogger) (*Data, error) {
	dat, err := os.ReadFile("book.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	text := string(dat)

	d := &Data{}
	d.Length = int64(len(text))
	d.Seconds = int64(time.Hour * 24 * 7)
	now := time.Now().UTC()
	d.SecondsPassed = ((int64(now.Nanosecond()) * int64(time.Nanosecond)) +
		(int64(now.Second()) * int64(time.Second)) +
		(int64(now.Minute()) * int64(time.Minute)) +
		(int64(now.Hour()) * int64(time.Hour)) +
		(int64(now.Weekday()) * int64(time.Hour*24)))

	d.Lookup = int64((float64(d.SecondsPassed) / float64(d.Seconds)) * float64(d.Length))
	if d.Lookup >= d.Length {
		d.Lookup = d.Length - 1
	}

	d.Character = rune(text[d.Lookup])

	l.Debugf("%s", d.Log())
	return d, nil
}

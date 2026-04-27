package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/icco/gutil/logging"
	"go.uber.org/zap"
)

var (
	log = logging.Must(logging.NewLogger("numbers"))
)

func main() {
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}
	log.Infow("Starting up", "host", fmt.Sprintf("http://localhost:%s", port))

	r := chi.NewRouter()
	r.Use(logging.Middleware(log.Desugar()))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		l := logging.FromContext(r.Context())
		d, err := GetCharacter(l)
		if err != nil {
			l.Errorw("get character", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%d", d.Character)
	})

	r.Get("/json", func(w http.ResponseWriter, r *http.Request) {
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
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("hi.")); err != nil {
			logging.FromContext(r.Context()).Errorw("write healthz", zap.Error(err))
		}
	})

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

// Data holds a single character lookup result derived from the current time.
type Data struct {
	Character     rune  `json:"character"`
	SecondsPassed int64 `json:"seconds_passed"`
	Seconds       int64 `json:"seconds"`
	Length        int64 `json:"length"`
	Lookup        int64 `json:"lookup"`
}

func (d *Data) Log() string {
	return fmt.Sprintf("(%v / %v) * %d = %d: %s (%d)", d.SecondsPassed, d.Seconds, d.Length, d.Lookup, string(d.Character), d.Character)
}

// GetCharacter returns the character in book.txt that corresponds to the
// current position within the week.
//
// The nanosecond component of now.Nanosecond() means SecondsPassed can
// exceed Seconds by up to 999_999_999 ns, making the raw Lookup index
// equal to (or greater than) Length.  The clamp below prevents the
// resulting index out-of-range panic.
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

	// Clamp: floating-point rounding can make Lookup == Length (or rarely
	// larger) during the final nanosecond of Saturday night UTC, causing a
	// panic on the slice access below.
	if d.Lookup >= d.Length {
		d.Lookup = d.Length - 1
	}

	d.Character = rune(text[d.Lookup])

	l.Debugf("%s", d.Log())
	return d, nil
}

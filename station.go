package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/icco/gutil/logging"
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
	r.Use(middleware.RealIP)
	r.Use(logging.Middleware(log.Desugar(), "icco-cloud"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		d, err := GetCharacter()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%d", d.Character)
	})

	r.Get("/json", func(w http.ResponseWriter, r *http.Request) {
		d, err := GetCharacter()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(d)
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi."))
	})

	log.Fatal(http.ListenAndServe(":"+port, r))
}

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

func GetCharacter() (*Data, error) {
	dat, err := os.ReadFile("book.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	d := &Data{}

	text := string(dat)
	d.Length = int64(len(text))
	d.Seconds = int64(time.Hour * 24 * 7)
	now := time.Now().UTC()
	d.SecondsPassed = ((int64(now.Nanosecond()) * int64(time.Nanosecond)) +
		(int64(now.Second()) * int64(time.Second)) +
		(int64(now.Minute()) * int64(time.Minute)) +
		(int64(now.Hour()) * int64(time.Hour)) +
		(int64(now.Weekday()) * int64(time.Hour*24)))

	d.Lookup = (d.SecondsPassed / d.Seconds) * d.Length
	d.Character = rune(text[d.Lookup])

	log.Debugf(d.Log())
	return d, nil
}

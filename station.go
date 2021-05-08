package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/icco/gutil/logging"
)

var (
	log = logging.Must(logging.NewLogger("numbers"))
)

func handler(w http.ResponseWriter, r *http.Request) {
	dat, err := ioutil.ReadFile("book.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	text := string(dat)
	length := len(text)
	seconds := float64(time.Hour * 24 * 7)
	now := time.Now().UTC()
	secPassed := ((int64(now.Nanosecond()) * int64(time.Nanosecond)) +
		(int64(now.Second()) * int64(time.Second)) +
		(int64(now.Minute()) * int64(time.Minute)) +
		(int64(now.Hour()) * int64(time.Hour)) +
		(int64(now.Weekday()) * int64(time.Hour*24)))

	lookup := int64((float64(secPassed) / seconds) * float64(length))
	char := rune(text[lookup])

	log.Debugf("(%v / %v) * %d = %d: %s (%d)", secPassed, seconds, length, lookup, string(char), char)
	fmt.Fprintf(w, "%d", char)
}

func main() {
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}
	log.Infow("Starting up", "host", fmt.Sprintf("http://localhost:%s", port))

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(logging.Middleware(log.Desugar(), "icco-cloud"))

	r.Get("/", handler)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi."))
	})

	log.Fatal(http.ListenAndServe(":"+port, r))
}

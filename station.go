package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dat, err := ioutil.ReadFile("book.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	length := len(dat)
	seconds := float64(604800)
	now := time.Now().UTC()
	seconds_passed := now.Second() + (now.Minute() * 60) + (now.Hour() * 3600) + (int(now.Weekday()) * 86400)
	lookup := int((float64(seconds_passed) / seconds) * float64(length))
	char := rune(dat[lookup])

	log.Printf("(%v / %v) * %d = %d: %s (%d)", seconds, seconds_passed, length, lookup, string(char), char)
	fmt.Fprintf(w, "%d", char)
}

func init() {
	http.HandleFunc("/", handler)
}

func main() {
	numbPtr := flag.Int("p", 8080, "Port to run on")
	flag.Parse()

	where := ":http"
	if *numbPtr != 80 {
		where = fmt.Sprintf(":%d", *numbPtr)
	}

	log.Fatal(http.ListenAndServe(where, nil))
}

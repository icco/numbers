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

	text := string(dat)
	length := len(text)
	seconds := float64(time.Hour * 24 * 7)
	now := time.Now().UTC()
	seconds_passed := ((int64(now.Nanosecond()) * int64(time.Nanosecond)) +
		(int64(now.Second()) * int64(time.Second)) +
		(int64(now.Minute()) * int64(time.Minute)) +
		(int64(now.Hour()) * int64(time.Hour)) +
		(int64(now.Weekday()) * int64(time.Hour*24)))

	lookup := int64((float64(seconds_passed) / seconds) * float64(length))
	char := rune(text[lookup])

	log.Printf("(%v / %v) * %d = %d: %s (%d)", seconds_passed, seconds, length, lookup, string(char), char)
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

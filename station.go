package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dat, err := ioutil.ReadFile("neuromancer.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	length := len(dat)
	seconds := float64(604800)
	now := time.Now().UTC()
	seconds_passed := now.Second() + (now.Minute() * 60) + (now.Hour() * 3600) + (int(now.Weekday()) * 86400)
	lookup := int((float64(seconds_passed) / seconds) * float64(length))

	fmt.Fprintf(w, "file is %d, seconds are %d of %d. %d is %s", length, seconds_passed, seconds, lookup, string(dat[lookup]))
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

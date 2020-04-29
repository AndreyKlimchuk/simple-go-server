package main

import (
    "log"
    "net/http"
    "html/template"
    "sync/atomic"
)

var indexTemplate = template.Must(template.ParseFiles("index.html"))
var visitorsCount uint64;

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		atomic.AddUint64(&visitorsCount, 1)
		err := indexTemplate.Execute(w, visitorsCount)
		if err != nil {
	        http.Error(w, err.Error(), http.StatusInternalServerError)
	    }
	} else {
		http.NotFound(w, r)
	}
}

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

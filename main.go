package main

import (
    "log"
    "net/http"
    "html/template"
    "sync/atomic"
    "flag"
    "strconv"
    "net"
)

var indexTemplate = template.Must(template.ParseFiles("index.html"))
var visitorsCount uint64;

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		atomic.AddUint64(&visitorsCount, 1)
		err := indexTemplate.Execute(w, visitorsCount)
		if err != nil {
			log.Println(err)
	    }
	} else {
		http.NotFound(w, r)
	}
}

func start_http_server(listener net.Listener) {
	http.HandleFunc("/", handler)
    log.Fatal(http.Serve(listener, nil))
}

func start_listener(port int) (net.Listener) {
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
	    log.Fatal(err)
	}
	return listener;
}

func main() {
	port := flag.Int("port", 8080, "port number")
	flag.Parse()
	start_http_server(start_listener(*port))
}

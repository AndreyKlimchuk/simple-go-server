package main

import (
	"testing"
	"net/http"
	"strconv"
	"sync"
	"os"
	"time"
	"net"
)

import "github.com/PuerkitoBio/goquery"

var port string;

func TestMain(m *testing.M) {
	listener := start_listener(0) // start listener on available port
	port = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	go start_http_server(listener)
	time.Sleep(1 * time.Second) // waiting server initialization
	os.Exit(m.Run())
}

func TestSequential(t *testing.T) {
	countOld := get_visitors_count(t)
	for i := 0; i < 10; i++ {
		countNew := get_visitors_count(t)
		if countOld + 1 != countNew {
			t.Fatalf("Invalid visitors count increment: %d -> %d", countOld, countNew)
		}
		countOld = countNew
	}
}

func TestParallel(t *testing.T) {
	var wg sync.WaitGroup

	// Trying to achive race condition when counter is not atomic
	countInit := get_visitors_count(t)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for c := 0; c < 100; c++ {
				dummy_request(t)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	countFinal := get_visitors_count(t)
	if countInit + 1001 != countFinal {
		t.Fatalf("Invalid visitors count increment: %d -> %d", countInit, countFinal)
	}
}

func request_page(t *testing.T) (*http.Response) {
	resp, err := http.Get("http://localhost:" + port +"/")
	if err != nil { 
		t.Fatal(err) 
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Bad response code: %d", resp.StatusCode) 
	}
	return resp  
}

func dummy_request(t *testing.T) {
	resp := request_page(t)
	resp.Body.Close()
}

func get_visitors_count(t *testing.T) (int) {
	resp := request_page(t) 
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	visitorsCount, err := strconv.Atoi(doc.Find("#visitors-count").Text())
	if err != nil {
		t.Fatal(err)
	}
	return visitorsCount
}
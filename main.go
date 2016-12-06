package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type History struct {
	sync.Mutex
	A, B, C EventPost
}

type EventPost struct {
	Timestamp int64
	Data      string
}

func main() {
	var port int
	flag.IntVar(&port, "port", 5555, "the port to which sendgrid will post")
	flag.Parse()

	hist := &History{}

	mux := http.NewServeMux()
	mux.HandleFunc("/event/api", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			data = []byte(err.Error())
		}
		hist.Lock()
		defer hist.Unlock()
		hist.C, hist.B, hist.A = hist.B, hist.A, EventPost{Data: string(data), Timestamp: time.Now().Unix()}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		hist.Lock()
		defer hist.Unlock()
		json.NewEncoder(w).Encode(hist)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatal(err)
	}
}

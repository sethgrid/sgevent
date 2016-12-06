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
	Data      []map[string]interface{}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 5555, "the port to which sendgrid will post")
	flag.Parse()

	hist := &History{}
	m := []map[string]interface{}{}

	mux := http.NewServeMux()
	mux.HandleFunc("/event/api", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			m = append(m, map[string]interface{}{"error": err.Error()})
		} else {
			err = json.Unmarshal(data, &m)
			if err != nil {
				m = append(m, map[string]interface{}{"error": err.Error()})
			}
		}
		hist.Lock()
		defer hist.Unlock()
		hist.C, hist.B, hist.A = hist.B, hist.A, EventPost{Data: m, Timestamp: time.Now().Unix()}
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

package heatmap

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type httpServer struct {
	storage datastore
	config  *config
}

type renderReturn struct {
	Target     string       `json:"target"`
	Datapoints []*datapoint `json:"datapoints"`
}

func (h *httpServer) version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("1.1.5\n"))
}

func (h *httpServer) functions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{}"))
}

func (h *httpServer) renderer(w http.ResponseWriter, r *http.Request) {
	targetArr, ok := r.URL.Query()["target"]
	target := ""

	if !ok {
		err := r.ParseForm()
		if err != nil {
			return
		}
		target = r.Form.Get("target")
	} else {
		target = targetArr[0]
	}

	from := time.Now()
	to := time.Now()

	w.Header().Set("Content-Type", "application/json")

	e := json.NewEncoder(w)
	rr := &renderReturn{Target: target, Datapoints: h.storage.Get(target, from, to)}
	e.Encode([]*renderReturn{rr})
}

// Main is the main entrypoint
func (h *httpServer) start() {
	http.HandleFunc("/functions", h.functions)
	http.HandleFunc("/version", h.version)
	http.HandleFunc("/render", h.renderer)
	log.Fatal(http.ListenAndServe(h.config.httpAddr, nil))
}

package heatmap

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type httpServer struct {
	storage datastore
	config  *config
}

type metricsFindReturn struct {
	Text          string `json:"text"`
	Expandable    int    `json:"expandable"`
	Leaf          int    `json:"leaf"`
	ID            string `json:"id"`
	AllowChildren int    `json:"allowChildren"`
}

type renderReturn struct {
	Target     string       `json:"target"`
	Datapoints []*datapoint `json:"datapoints"`
}

func (h *httpServer) tags(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func (h *httpServer) version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("1.1.5\n"))
}

func (h *httpServer) functions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{}"))
}

func lastSegment(key string) string {
	i := strings.LastIndex(key, ".")
	if i != -1 && i+1 < len(key) {
		return key[i+1:]
	}
	return key
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func (h *httpServer) metricsFind(w http.ResponseWriter, r *http.Request) {
	queryArr, ok := r.URL.Query()["query"]
	query := ""

	if !ok {
		err := r.ParseForm()
		if err != nil {
			return
		}
		query = r.Form.Get("query")
	} else {
		query = queryArr[0]
	}

	results := []*metricsFindReturn{}

	for _, globResult := range h.storage.Glob(query) {
		results = append(results, &metricsFindReturn{
			Text:          lastSegment(globResult.name),
			Expandable:    boolToInt(globResult.hasChildren),
			Leaf:          boolToInt(globResult.isLeaf),
			ID:            globResult.name,
			AllowChildren: boolToInt(globResult.hasChildren),
		})
	}

	e := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	e.Encode(results)
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
	resultArray := []*renderReturn{}

	globbedTargets := h.storage.Glob(target)
	if len(globbedTargets) > 0 {
		globbedName := globbedTargets[0].name
		resultArray = append(resultArray, &renderReturn{
			Target:     globbedName,
			Datapoints: h.storage.Get(globbedName, from, to),
		})
	}

	e := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	e.Encode(resultArray)
}

// Main is the main entrypoint
func (h *httpServer) start() {
	http.HandleFunc("/tags/autoComplete/tags", h.tags)
	http.HandleFunc("/metrics/find", h.metricsFind)
	http.HandleFunc("/functions", h.functions)
	http.HandleFunc("/version", h.version)
	http.HandleFunc("/render", h.renderer)
	log.Fatal(http.ListenAndServe(h.config.httpAddr, nil))
}

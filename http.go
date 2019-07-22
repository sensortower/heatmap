package heatmap

import (
	"encoding/json"
	"net/http"
	"sort"
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

func requestParam(r *http.Request, name string) (val string) {
	arr, ok := r.URL.Query()[name]

	if !ok {
		err := r.ParseForm()
		if err != nil {
			return
		}
		val = r.Form.Get(name)
	} else {
		val = arr[0]
	}
	return
}

type sortableMetricsFindReturn []*metricsFindReturn

func (m sortableMetricsFindReturn) Len() int {
	return len(m)
}

func (m sortableMetricsFindReturn) Less(i, j int) bool {
	return strings.Compare(m[i].Text, m[j].Text) < 0
}

func (m sortableMetricsFindReturn) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (h *httpServer) metricsFind(w http.ResponseWriter, r *http.Request) {
	results := []*metricsFindReturn{}

	for _, globResult := range h.storage.Glob(requestParam(r, "query")) {
		results = append(results, &metricsFindReturn{
			Text:          lastSegment(globResult.name),
			Expandable:    boolToInt(globResult.hasChildren),
			Leaf:          boolToInt(globResult.isLeaf),
			ID:            globResult.name,
			AllowChildren: boolToInt(globResult.hasChildren),
		})
	}

	sort.Sort(sortableMetricsFindReturn(results))

	e := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	e.Encode(results)
}

func (h *httpServer) renderer(w http.ResponseWriter, r *http.Request) {
	// TODO: use the ones from request
	from := time.Now()
	to := time.Now()
	resultArray := []*renderReturn{}

	globbedTargets := h.storage.Glob(requestParam(r, "target"))
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
	logError.Fatalln(http.ListenAndServe(h.config.httpAddr, nil))
}

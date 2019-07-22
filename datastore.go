package heatmap

import (
	"encoding/json"
	"time"
)

type datapoint struct {
	timestamp uint32
	value     float32
}

func (dp *datapoint) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{dp.value, dp.timestamp})
}

type globResult struct {
	name        string
	isLeaf      bool
	hasChildren bool
}

type datastore interface {
	Glob(key string) []*globResult
	Get(key string, from, to time.Time) []*datapoint
	Put(key string, p *datapoint)
}

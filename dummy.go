package heatmap

import (
	"math"
	"math/rand"
	"time"
)

type dummyData struct {
	storage datastore
	config  *config
}

func (dd *dummyData) start() {
	t := time.NewTicker(time.Second)
	for {
		<-t.C
		ts := time.Now()
		for i := 0; i < rand.Intn(100); i++ {
			dd.storage.Put("dummy-data", &datapoint{timestamp: ts, duration: math.Abs(rand.NormFloat64() * 1000)})
		}
	}
}

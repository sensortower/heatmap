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
		ts := uint32(time.Now().Unix())
		for i := 0; i < rand.Intn(100); i++ {
			dd.storage.Put("dummy-data", &datapoint{timestamp: ts, value: float32(math.Abs(rand.NormFloat64() * 1000))})
		}
	}
}

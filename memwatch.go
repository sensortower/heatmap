package heatmap

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type memWatch struct {
	ramDatastore *ramDatastore
	config       *config

	memThreshold uint64 // in bytes
	gcChan       chan struct{}
}

var multipliers = map[string]uint64{
	"bytes": 1,
	"b":     1,
	"kb":    1024,
	"mb":    1024 * 1024,
	"gb":    1024 * 1024 * 1024,
	"tb":    1024 * 1024 * 1024 * 1024,
	"pb":    1024 * 1024 * 1024 * 1024 * 1024,
}

func (m *memWatch) scheduleGCRun() {
	logDebug.Printf("[MEMWATCH] scheduling a GC run")
	select {
	case m.gcChan <- struct{}{}:
	default:
	}
}

func (m *memWatch) gcSubroutine() {
	for {
		<-m.gcChan
		reportTime("GC run", func() { runtime.GC() })
	}
}

func (m *memWatch) start() {
	m.gcChan = make(chan struct{}, 1)
	go m.gcSubroutine()

	t := time.NewTicker(time.Second)

	if strings.Contains(m.config.memThreshold, "%") {
		percentage, err := strconv.ParseFloat(strings.Replace(m.config.memThreshold, "%", "", 1), 64)
		if err != nil {
			panic(err)
		}

		m.memThreshold = uint64(percentage / 100.0 * float64(memoryTotal()))
	} else {
		re := regexp.MustCompile("(\\d*\\.?\\d*)\\s*(\\S+)")
		submatches := re.FindStringSubmatch(m.config.memThreshold)
		if len(submatches) < 3 {
			panic(fmt.Sprintf("could not parse mem-threshold %s", m.config.memThreshold))
		}
		value, err := strconv.ParseFloat(submatches[1], 64)
		if err != nil {
			panic(err)
		}
		multiplier := strings.ToLower(submatches[2])
		if multiplierUint, ok := multipliers[multiplier]; ok {
			m.memThreshold = uint64(value * float64(multiplierUint))
		} else {
			panic(fmt.Sprintf("could not parse mem-threshold, unknown multipler %s", multiplier))
		}
	}

	for {
		<-t.C
		used := memoryUsed()
		logDebug.Printf("[MEMWATCH] checking on memory usage %d/%d", used, m.memThreshold)
		if used > m.memThreshold {
			m.scheduleGCRun()
		}

		time.Sleep(time.Second)

		used = memoryUsed()
		if used > m.memThreshold {
			logDebug.Printf("[MEMWATCH] doing a full cleanup")
			reportTime("RAM datastore cleanup", func() { m.ramDatastore.cleanup() })
			m.scheduleGCRun()
		}
	}
}

func memoryUsed() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func reportTime(label string, cb func()) {
	t := time.Now()
	cb()
	logDebug.Printf("[MEMWATCH] %s took %s", label, time.Now().Sub(t).String())
}

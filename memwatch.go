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

func (m *memWatch) start() {
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

		i := 0
		logDebug.Printf("[MEMWATCH] checking on memory usage %d/%d", used, m.memThreshold)
		for i < 3 && used > m.memThreshold {
			m.ramDatastore.cleanup()
			runtime.GC()
			i++
		}
	}
}

func memoryUsed() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Sys
}

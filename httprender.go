package heatmap

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const minXBucketSize = uint32(5)
const maxXBucketSize = uint32(3600)

func minMax(val, min, max uint32) uint32 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(422)
	w.Write([]byte(fmt.Sprintf("error parsing request: %v", err)))
}

func (h *httpServer) renderer(w http.ResponseWriter, r *http.Request) {
	from := parseATTime(requestParam(r, "from"))
	to := parseATTime(requestParam(r, "to"))

	maxDataPoints, err := strconv.Atoi(requestParam(r, "maxDataPoints"))
	if err != nil {
		maxDataPoints = 500
	}

	maxDataPoints /= 10

	resultArray := []*renderReturn{}

	target := requestParam(r, "target")
	inBuckets := strings.HasPrefix(target, "inBuckets")
	target = strings.Replace(target, "inBuckets", "", 1)
	target = strings.Trim(target, "()")
	targetArgs := strings.Split(target, ",")
	for i, t := range targetArgs {
		targetArgs[i] = strings.Trim(strings.TrimSpace(t), "\"'")
	}
	target = targetArgs[0]

	bucketCount := 10
	if len(targetArgs) > 1 && targetArgs[1] != "auto" {
		bucketCount, err = strconv.Atoi(targetArgs[1])
		if err != nil {
			handleError(w, err)
			return
		}
	}

	yLimit := float32(0.0)
	if len(targetArgs) > 2 && targetArgs[2] != "auto" {
		l64, err := strconv.ParseFloat(targetArgs[2], 64)
		if err != nil {
			handleError(w, err)
			return
		}
		yLimit = float32(l64)
	}

	logScale := false
	if len(targetArgs) > 3 && targetArgs[3] != "auto" {
		logScale = targetArgs[3] == "true"
	}

	globbedTargets := h.storage.Glob(target)
	if len(globbedTargets) > 0 {
		globbedName := globbedTargets[0].name
		allData := h.storage.Get(globbedName, from, to)

		if inBuckets {
			minYVal := float32(0.0)
			maxYVal := yLimit
			if yLimit == 0.0 {
				maxYVal = 1.0
				for _, d := range allData {
					if d.value > maxYVal {
						maxYVal = d.value
					}
				}
			}

			if logScale {
				maxYVal = float32(math.Log(float64(maxYVal)))
			}

			bucketXSize := uint32(to.Sub(from)/time.Second) / uint32(maxDataPoints)
			bucketXSize = minMax(bucketXSize, minXBucketSize, maxXBucketSize)
			bucketYSize := (maxYVal - minYVal) / float32(bucketCount-1)

			bucketsMap := make([]map[uint32]*datapoint, bucketCount)
			for i := range bucketsMap {
				bucketsMap[i] = make(map[uint32]*datapoint)
			}
			for _, d := range allData {
				dv := d.value
				if logScale {
					dv = float32(math.Log(float64(dv)))
				}
				bucketIndex := int(dv / bucketYSize)
				if bucketIndex < bucketCount {
					key := d.timestamp / bucketXSize * bucketXSize
					if v, ok := bucketsMap[bucketIndex][key]; ok {
						v.value += 1.0
					} else {
						bucketsMap[bucketIndex][key] = &datapoint{timestamp: key, value: 1.0}
					}
				}
			}
			for i, b := range bucketsMap {
				datapoints := make([]*datapoint, 0, len(b))

				for _, value := range b {
					datapoints = append(datapoints, value)
				}
				bucketLowerBoundary := float32(i) * bucketYSize
				if logScale {
					bucketLowerBoundary = float32(math.Exp(float64(bucketLowerBoundary)))
				}
				resultArray = append(resultArray, &renderReturn{
					Target:     fmt.Sprintf("%f", bucketLowerBoundary),
					Datapoints: datapoints,
				})
			}
		} else {
			resultArray = append(resultArray, &renderReturn{
				Target:     globbedName,
				Datapoints: allData,
			})
		}
	}

	e := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	e.Encode(resultArray)
}

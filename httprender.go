package heatmap

import (
	"encoding/json"
	"fmt"
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
	target = strings.Trim(target, "\"'")

	globbedTargets := h.storage.Glob(target)
	if len(globbedTargets) > 0 {
		globbedName := globbedTargets[0].name
		allData := h.storage.Get(globbedName, from, to)

		if inBuckets {
			bucketCount := 10
			minVal := float32(0.0)
			maxVal := float32(1.0)
			for _, d := range allData {
				if d.duration > maxVal {
					maxVal = d.duration
				}
			}

			bucketXSize := uint32(to.Sub(from)/time.Second) / uint32(maxDataPoints)
			bucketXSize = minMax(bucketXSize, minXBucketSize, maxXBucketSize)
			bucketYSize := (maxVal - minVal) / float32(bucketCount-1)

			bucketsMap := make([]map[uint32]*datapoint, bucketCount)
			for i := range bucketsMap {
				bucketsMap[i] = make(map[uint32]*datapoint)
			}
			for _, d := range allData {
				bucketIndex := int(d.duration / bucketYSize)
				key := d.timestamp / bucketXSize * bucketXSize
				if v, ok := bucketsMap[bucketIndex][key]; ok {
					v.duration += 1.0
				} else {
					bucketsMap[bucketIndex][key] = &datapoint{timestamp: key, duration: 1.0}
				}
			}
			for i, b := range bucketsMap {
				datapoints := make([]*datapoint, 0, len(b))

				for _, value := range b {
					datapoints = append(datapoints, value)
				}
				resultArray = append(resultArray, &renderReturn{
					Target:     fmt.Sprintf("%f", float32(i)*bucketYSize),
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

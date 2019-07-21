# heatmap

`heatmap` is an in-memory time series database with statsd interface for incoming data and graphite-compatible rendering API.

It's intended to be used with Grafana's [Heatmap](https://grafana.com/docs/features/panels/heatmap/) feature. It's inspired by [Brendan Gregg's HeatMap](http://www.brendangregg.com/HeatMaps/latency.html) project.


Install heatmap
```
go get github.com/sensortower/heatmap/cmd/heatmap
```

Usage of heatmap:
```
  -create-dummy-data
    	flag that enables dummy data generation
  -http-addr string
    	http endpoint address to bind to (default ":10000")
  -statsd-addr string
    	statsd address to bind to (default ":8125")

```

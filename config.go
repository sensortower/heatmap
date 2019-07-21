package heatmap

import "flag"

type config struct {
	logLevel        string
	statsdAddr      string
	httpAddr        string
	memThreshold    string
	createDummyData bool
}

func (c *config) populateFromFlags() {
	flag.StringVar(&c.logLevel, "log-level", "info", "log level (debug / info / error)")
	flag.StringVar(&c.statsdAddr, "statsd-addr", ":8125", "statsd address to bind to")
	flag.StringVar(&c.httpAddr, "http-addr", ":10000", "http endpoint address to bind to")
	flag.StringVar(&c.memThreshold, "mem-threshold", "80%", "maximum amount of memory heatmap is allowed to take")
	flag.BoolVar(&c.createDummyData, "create-dummy-data", false, "flag that enables dummy data generation")
	flag.Parse()
}

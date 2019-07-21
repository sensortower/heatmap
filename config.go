package heatmap

import "flag"

type config struct {
	statsdAddr      string
	httpAddr        string
	createDummyData bool
}

func (c *config) populateFromFlags() {
	flag.StringVar(&c.statsdAddr, "statsd-addr", ":8125", "statsd address to bind to")
	flag.StringVar(&c.httpAddr, "http-addr", ":10000", "http endpoint address to bind to")
	flag.BoolVar(&c.createDummyData, "create-dummy-data", false, "flag that enables dummy data generation")
	flag.Parse()
}

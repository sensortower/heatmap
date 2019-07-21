package heatmap

// Main is the main entrypoint
func Main() {
	config := &config{}
	config.populateFromFlags()

	storage := newRAMDatastore()

	l := &statsdUDPListener{storage, config}
	s := &httpServer{storage, config}
	go l.start()
	go s.start()

	if config.createDummyData {
		d := &dummyData{storage, config}
		go d.start()
	}

	select {}
}

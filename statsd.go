package heatmap

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"
)

type message struct {
	key       string
	datapoint *datapoint
}

type statsdUDPListener struct {
	storage datastore
	config  *config
}

// # key:value|type
func parseAndFilter(msg []byte) (*message, error) {
	i := bytes.Index(msg, []byte{'|'})
	if i == -1 {
		return nil, fmt.Errorf("message is invalid: %q", msg)
	}

	// # "abc|" i = 3, len = 4
	if i+1 >= len(msg) {
		return nil, fmt.Errorf("message is invalid: %q", msg)
	}

	msgType := msg[i+1]

	if msgType != 'm' {
		return nil, nil
	}

	i2 := bytes.Index(msg, []byte{':'})

	// # "a|c:" i = 3, len = 4
	if i2 > i {
		return nil, fmt.Errorf("message is invalid: %q", msg)
	}

	val, err := strconv.ParseFloat(string(msg[i2+1:i]), 32)

	if err != nil {
		return nil, fmt.Errorf("message is invalid: %q", msg)
	}

	d := &datapoint{
		timestamp: uint32(time.Now().Unix()),
		value:     float32(val),
	}

	m := &message{
		key:       string(msg[:i2]),
		datapoint: d,
	}

	return m, nil
}

func (sl *statsdUDPListener) start() error {
	conn, err := net.ListenPacket("udp", sl.config.statsdAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := make([]byte, 4096)
	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			return err
		}
		newBuf := buf[:n]
		slices := bytes.Split(newBuf, []byte("\n"))
		for _, slice := range slices {
			if len(slice) == 0 {
				continue
			}
			m, err := parseAndFilter(slice)
			if err != nil {
				logError.Println("[STATSD] failed to parse message:", err)
			} else if m != nil {
				sl.storage.Put(m.key, m.datapoint)
			}
		}
	}
}

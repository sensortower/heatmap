// +build darwin

package heatmap

import (
	"encoding/binary"
	"syscall"
)

func memoryTotal() uint64 {
	str, err := syscall.Sysctl("hw.memsize")
	if err != nil {
		panic(err)
	}

	b := []byte(str)
	b = append(b, 0)
	val := binary.LittleEndian.Uint64(b)

	return val
}

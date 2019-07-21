// +build linux

package heatmap

import "syscall"

func memoryTotal() uint64 {
	st := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(st)
	if err != nil {
		panic(err)
	}
	return uint64(st.Totalram) * uint64(st.Unit)
}

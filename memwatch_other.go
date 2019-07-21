// +build !linux,!darwin

package heatmap

func memoryTotal() uint64 {
	panic("not implemented on this platform")
	return 0
}

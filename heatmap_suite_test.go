package heatmap

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHeatmap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Heatmap Suite")
}

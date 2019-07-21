package heatmap

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ramDatastore", func() {
	Context("Put", func() {
		It("works with simple cases", func() {
			r := newRAMDatastore()
			dp := &datapoint{
				timestamp: time.Now(),
				duration:  0.33,
			}
			r.Put("test", dp)

			arr := r.root.children["test"].data
			Expect(len(arr)).To(Equal(1))
			Expect(arr[0]).To(Equal(dp))
		})

		It("works with slightly more complicated cases", func() {
			r := newRAMDatastore()
			dp := &datapoint{
				timestamp: time.Now(),
				duration:  0.33,
			}
			r.Put("foo.bar.baz", dp)

			arr := r.root.children["foo"].children["bar"].children["baz"].data
			Expect(len(arr)).To(Equal(1))
			Expect(arr[0]).To(Equal(dp))
		})
	})

	Context("Get", func() {
		It("works with simple cases", func() {
			r := newRAMDatastore()
			dp := &datapoint{
				timestamp: time.Now(),
				duration:  0.33,
			}
			r.Put("test", dp)

			arr := r.Get("test", time.Now(), time.Now())
			Expect(len(arr)).To(Equal(1))
			Expect(arr[0]).To(Equal(dp))
		})

		It("works with slightly more complicated cases", func() {
			r := newRAMDatastore()
			dp := &datapoint{
				timestamp: time.Now(),
				duration:  0.33,
			}
			r.Put("foo.bar.baz", dp)

			arr := r.Get("foo.bar.baz", time.Now(), time.Now())
			Expect(len(arr)).To(Equal(1))
			Expect(arr[0]).To(Equal(dp))
		})
	})

	Context("Glob", func() {
		It("works with simple cases", func() {
			r := newRAMDatastore()
			dp := &datapoint{
				timestamp: time.Now(),
				duration:  0.33,
			}
			r.Put("prefix.foo.suffix", dp)
			r.Put("prefix.bar.suffix", dp)
			r.Put("prefix.baz.suffix", dp)

			arr := r.Glob("prefix.*.suffix")
			Expect(arr).To(ConsistOf([]*globResult{
				&globResult{
					name:   "prefix.foo.suffix",
					isLeaf: true,
				}, &globResult{
					name:   "prefix.bar.suffix",
					isLeaf: true,
				}, &globResult{
					name:   "prefix.baz.suffix",
					isLeaf: true,
				}}))
		})

		It("works with asterisk being the last element", func() {
			r := newRAMDatastore()
			dp := &datapoint{
				timestamp: time.Now(),
				duration:  0.33,
			}
			r.Put("prefix.foo.suffix", dp)
			r.Put("prefix.bar.suffix", dp)
			r.Put("prefix.baz.suffix", dp)

			arr := r.Glob("prefix.*")
			Expect(arr).To(ConsistOf([]*globResult{
				&globResult{
					name:        "prefix.foo",
					hasChildren: true,
				}, &globResult{
					name:        "prefix.bar",
					hasChildren: true,
				}, &globResult{
					name:        "prefix.baz",
					hasChildren: true,
				},
			}))
		})

		It("works with a nonexisting root element", func() {
			r := newRAMDatastore()
			dp := &datapoint{
				timestamp: time.Now(),
				duration:  0.33,
			}
			r.Put("prefix.foo.suffix", dp)
			r.Put("prefix.bar.suffix", dp)
			r.Put("prefix.baz.suffix", dp)

			arr := r.Glob("404.*")
			Expect(arr).To(BeEmpty())
		})

		It("works with a nonexisting children element", func() {
			r := newRAMDatastore()
			dp := &datapoint{
				timestamp: time.Now(),
				duration:  0.33,
			}
			r.Put("prefix.foo.suffix", dp)
			r.Put("prefix.bar.suffix", dp)
			r.Put("prefix.baz.suffix", dp)

			arr := r.Glob("prefix.404")
			Expect(arr).To(BeEmpty())
		})

		It("works with root asterisk case", func() {
			r := newRAMDatastore()
			dp := &datapoint{
				timestamp: time.Now(),
				duration:  0.33,
			}
			r.Put("foo", dp)
			r.Put("bar", dp)
			r.Put("baz", dp)

			arr := r.Glob("*")
			Expect(arr).To(ConsistOf([]*globResult{
				&globResult{
					name:   "foo",
					isLeaf: true,
				}, &globResult{
					name:   "bar",
					isLeaf: true,
				}, &globResult{
					name:   "baz",
					isLeaf: true,
				},
			}))
		})
	})
})

package heatmap

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type datapoint struct {
	timestamp time.Time
	duration  float64
}

func (dp *datapoint) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{dp.duration, dp.timestamp.Unix()})
}

type datastore interface {
	Glob(key string) []string
	Get(key string, from, to time.Time) []*datapoint
	Put(key string, p *datapoint)
}

type treeNode struct {
	name          string
	parent        *treeNode
	childrenMutex sync.Mutex
	children      map[string]*treeNode

	data []*datapoint
}

func (tn *treeNode) glob(res *[]string, prefix string, fragments []string) {
	fmt.Println(prefix, fragments, res)
	for i, f := range fragments {
		if prefix != "" {
			prefix += "."
		}
		if tn == nil {
			break
		}
		if f == "*" {
			for k, child := range tn.children {
				child.glob(res, prefix+k, fragments[i+1:])
			}
			return
		}
		prefix += f
		tn = tn.children[f]
	}
	*res = append(*res, prefix)
	return
}

type ramDatastore struct {
	root *treeNode
}

func newRAMDatastore() *ramDatastore {
	return &ramDatastore{
		root: &treeNode{
			name:     "",
			parent:   nil,
			children: make(map[string]*treeNode),
		},
	}
}

func (rd *ramDatastore) Put(key string, p *datapoint) {
	fragments := strings.Split(key, ".")
	tn := rd.root
	for _, f := range fragments {
		if n, ok := tn.children[f]; ok {
			tn = n
		} else {
			newNode := &treeNode{name: f, parent: tn, children: make(map[string]*treeNode)}
			tn.children[f] = newNode
			tn = newNode
		}
	}
	tn.data = append(tn.data, p)
}

func (rd *ramDatastore) Get(key string, from, to time.Time) (res []*datapoint) {
	fragments := strings.Split(key, ".")
	tn := rd.root
	for _, f := range fragments {
		if n, ok := tn.children[f]; ok {
			tn = n
		} else {
			return
		}
	}

	// TODO: add time limiting
	res = append(res, tn.data...)

	return
}

func (rd *ramDatastore) Glob(key string) (res []string) {
	fragments := strings.Split(key, ".")
	rd.root.glob(&res, "", fragments)
	return
}

package heatmap

import (
	"strings"
	"sync"
	"time"
)

type treeNode struct {
	name          string
	parent        *treeNode
	childrenMutex sync.Mutex
	children      map[string]*treeNode

	data []*datapoint
}

func (tn *treeNode) glob(res *[]*globResult, prefix string, fragments []string) {
	for i, f := range fragments {
		if prefix != "" {
			prefix += "."
		}
		if tn == nil {
			return
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
	*res = append(*res, &globResult{
		name:        prefix,
		isLeaf:      len(tn.data) > 0,
		hasChildren: len(tn.children) > 0,
	})
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

func (rd *ramDatastore) Glob(key string) (res []*globResult) {
	fragments := strings.Split(key, ".")
	rd.root.glob(&res, "", fragments)
	return
}

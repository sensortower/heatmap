package heatmap

import (
	"regexp"
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

func globPatternToRegexp(pattern string) *regexp.Regexp {
	expr := "^"

	var multipleChoice bool

	for _, ch := range pattern {
		if ch == '*' {
			expr += ".*?"
		} else if ch == '{' {
			expr += "("
			multipleChoice = true
		} else if ch == '}' {
			multipleChoice = false
			expr += ")"
		} else if multipleChoice && ch == ',' {
			expr += "|"
		} else {
			expr += string(ch)
		}
	}

	expr += "$"
	return regexp.MustCompile(expr)
}

func (tn *treeNode) recursiveCleanup() {
	tn.data = tn.data[len(tn.data)/10:]
	for _, child := range tn.children {
		child.recursiveCleanup()
	}
}

func (tn *treeNode) glob(res *[]*globResult, prefix string, fragments []string) {
	if len(fragments) == 0 {
		if tn == nil {
			return
		}
		*res = append(*res, &globResult{
			name:        prefix,
			isLeaf:      len(tn.data) > 0,
			hasChildren: len(tn.children) > 0,
		})
		return
	}

	if prefix != "" {
		prefix += "."
	}

	r := globPatternToRegexp(fragments[0])
	for k, child := range tn.children {
		if r.MatchString(k) {
			child.glob(res, prefix+k, fragments[1:])
		}
	}
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

	fromInt := uint32(from.Unix())
	toInt := uint32(to.Unix())
	for _, d := range tn.data {
		if d.timestamp >= fromInt && d.timestamp <= toInt {
			res = append(res, d)
		}
	}

	return
}

func (rd *ramDatastore) Glob(key string) (res []*globResult) {
	fragments := strings.Split(key, ".")
	rd.root.glob(&res, "", fragments)
	return
}

func (rd *ramDatastore) cleanup() {
	rd.root.recursiveCleanup()
}

package ahocorasick

import (
	"github.com/koron/nvcheck/internal/trie"
)

type Event struct {
	root *trie.TernaryNode
	curr *trie.TernaryNode
}

func (ev *Event) Next() *Data {
	for ev.curr != ev.root {
		d := getNodeData(ev.curr)
		ev.curr = d.failure
		if d.Pattern != nil {
			return d
		}
	}
	return nil
}

package ahocorasick

import (
	"github.com/koron/nvcheck/internal/trie"
)

// Event represents matching posibility events.
type Event struct {
	root *trie.TernaryNode
	curr *trie.TernaryNode
}

// Next returns details of match event.
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

package ahocorasick

import (
	"github.com/koron/nvcheck/internal/trie"
)

// Iter represents iteratable matcher (state machine).
type Iter struct {
	trie *trie.TernaryTrie
	root *trie.TernaryNode
	curr *trie.TernaryNode
}

// Put puts a rune to state machine and gets Event.
func (it *Iter) Put(r rune) *Event {
	it.curr = getNextNode(it.curr, it.root, r)
	if it.curr == it.root {
		return nil
	}
	return &Event{
		root: it.root,
		curr: it.curr,
	}
}

// Has return true when current state has node/path for rune r.
func (it *Iter) Has(r rune) bool {
	n := it.curr.Get(r)
	return n != nil
}

// Reset resets state machine to start state.
func (it *Iter) Reset() {
	it.curr = it.root
}

package ahocorasick

import (
	"github.com/koron/nvcheck/internal/trie"
)

type Iter struct {
	trie *trie.TernaryTrie
	root *trie.TernaryNode
	curr *trie.TernaryNode
}

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

func (it *Iter) Reset() {
	it.curr = it.root
}

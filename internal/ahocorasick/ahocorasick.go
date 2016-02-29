package ahocorasick

import (
	"unicode/utf8"

	"github.com/koron/nvcheck/internal/trie"
)

// Matcher is strings matcher uses aho-corasick algorithm.
type Matcher struct {
	trie *trie.TernaryTrie
}

// Match represents a match by Matcher.
type Match struct {
	Index   int
	Pattern string
	Value   interface{}
}

// Data represents details of match.
type Data struct {
	Pattern *string
	Offset  int
	Value   interface{}

	failure *trie.TernaryNode
}

// New creates a new matcher.
func New() *Matcher {
	return &Matcher{
		trie: trie.NewTernaryTrie(),
	}
}

// Add adds a string pattern with user value v.
func (m *Matcher) Add(pattern string, v interface{}) {
	_, n := utf8.DecodeLastRuneInString(pattern)
	m.trie.Put(pattern, &Data{
		Pattern: &pattern,
		Offset:  len(pattern) - n,
		Value:   v,
	})
}

// Compile compiles a matcher for matching.
func (m *Matcher) Compile() error {
	m.trie.Balance()
	root := m.trie.Root().(*trie.TernaryNode)
	root.SetValue(&Data{failure: root})
	// fill data.failure of each node.
	trie.EachWidth(m.trie, func(n trie.Node) bool {
		parent := n.(*trie.TernaryNode)
		parent.Each(func(m trie.Node) bool {
			fillFailure(m.(*trie.TernaryNode), root, parent)
			return true
		})
		return true
	})
	return nil
}

// Iter creates an Iter instance.
func (m *Matcher) Iter() *Iter {
	r := m.trie.Root().(*trie.TernaryNode)
	return &Iter{
		trie: m.trie,
		root: r,
		curr: r,
	}
}

func fillFailure(curr, root, parent *trie.TernaryNode) {
	data := getNodeData(curr)
	if data == nil {
		data = &Data{}
		curr.SetValue(data)
	}
	if parent == root {
		data.failure = root
		return
	}
	// Determine failure node.
	fnode := getNextNode(getNodeFailure(parent, root), root, curr.Label())
	data.failure = fnode
}

// Match returns a channel to stream all Matches.
func (m *Matcher) Match(text string) <-chan Match {
	ch := make(chan Match, 1)
	go m.startMatch(text, ch)
	return ch
}

func (m *Matcher) startMatch(text string, ch chan<- Match) {
	defer close(ch)
	root := m.trie.Root().(*trie.TernaryNode)
	curr := root
	for i, r := range text {
		curr = getNextNode(curr, root, r)
		if curr == root {
			continue
		}
		fireAll(curr, root, ch, i)
	}
}

func getNextNode(node, root *trie.TernaryNode, r rune) *trie.TernaryNode {
	for {
		next, _ := node.Get(r).(*trie.TernaryNode)
		if next != nil {
			return next
		} else if node == root {
			return root
		}
		node = getNodeFailure(node, root)
	}
}

func fireAll(curr, root *trie.TernaryNode, ch chan<- Match, idx int) {
	for curr != root {
		data := getNodeData(curr)
		if data.Pattern != nil {
			ch <- Match{
				Index:   idx - data.Offset,
				Pattern: *data.Pattern,
				Value:   data.Value,
			}
		}
		curr = data.failure
	}
}

func getNodeData(node *trie.TernaryNode) *Data {
	d, _ := node.Value().(*Data)
	return d
}

func getNodeFailure(node, root *trie.TernaryNode) *trie.TernaryNode {
	next := getNodeData(node).failure
	if next == nil {
		return root
	}
	return next
}

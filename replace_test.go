package main

import (
	"bytes"
	"testing"
)

func testReplace(t *testing.T, d Dict, name, src, expected string) {
	m, err := d.toM()
	if err != nil {
		t.Fatal(err)
	}
	c := &ctx{m: m, fname: name}
	c.setup(src)
	if err := c.find(); err != nil {
		t.Errorf("%s: unexpected error: %s", name, err)
		return
	}
	buf := new(bytes.Buffer)
	replaceFounds(c, buf)
	actual := buf.String()
	if actual != expected {
		t.Errorf("%s: not match: expected=%q actual=%q",
			name, expected, actual)
	}
}

func TestReplace(t *testing.T) {
	d := parseDict(t, dict0)

	testReplace(t, d, "noop", "ユーザー", "ユーザー")
	testReplace(t, d, "simple", "ユーザ", "ユーザー")
	testReplace(t, d, "keep white", "ユー\nザ", "ユー\nザー")
	testReplace(t, d, "keep more whites", "  ユー\n  ザ", "  ユー\n  ザー")

	testReplace(t, d, "suppress control", "なに", "何")
	testReplace(t, d, "suppress", "どんなに", "どんなに")

	testReplace(t, d, "word includes white", "foobar", "foo bar")
	testReplace(t, d, "word includes white 1", "foo bar", "foo bar")
	testReplace(t, d, "word includes white 2", "foo  bar", "foo  bar")
	testReplace(t, d, "word includes white 3", "foo\nbar", "foo\nbar")
}

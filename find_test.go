package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type errs struct {
	v []error
}

func (e *errs) put(err error) {
	if err == nil {
		return
	}
	e.v = append(e.v, err)
}

func (e *errs) err() error {
	if len(e.v) == 0 {
		return nil
	}
	var b bytes.Buffer
	for i, err := range e.v {
		if i != 0 {
			b.WriteString(". ")
		}
		b.WriteString(err.Error())
	}
	return errors.New(b.String())
}

type found struct {
	begin int
	end   int
	text  string
	fix   string
}

func (f *found) match(t *Found) error {
	var errs errs
	if f.begin != t.Begin {
		errs.put(fmt.Errorf("begin expected %d but actulaly %d",
			f.begin, t.Begin))
	}
	if f.end != t.End {
		errs.put(fmt.Errorf("end expected %d but actulaly %d",
			f.end, t.End))
	}
	if f.text != t.Word.Text {
		errs.put(fmt.Errorf("text expected %q but actulaly %q",
			f.text, t.Word.Text))
	}
	if f.fix != "" {
		if t.Word.Fix == nil {
			errs.put(fmt.Errorf("less fix: %q", f.fix))
		} else if f.fix != *t.Word.Fix {
			errs.put(fmt.Errorf("fix expected %q but actulaly %q",
				f.fix, *t.Word.Fix))
		}
	} else if t.Word.Fix != nil {
		errs.put(fmt.Errorf("much fix: %q", *t.Word.Fix))
	}
	return errs.err()
}

func testFind(t *testing.T, d Dict, name string, expected []found, s string) {
	m, err := d.toM()
	if err != nil {
		t.Fatal(err)
	}
	c := &ctx{m: m, fname: name}
	c.setup(s)
	if err := c.find(); err != nil {
		t.Errorf("%s: unexpected error: %s", name, err)
		return
	}
	// check founds against with expected.
	for i, f := range expected {
		if i >= len(c.founds) {
			t.Errorf("%s: less founds: expected %d actually %d",
				name, len(expected), len(c.founds))
			return
		}
		err := f.match(c.founds[i])
		if err != nil {
			t.Errorf("%s: not match at %d: %s", name, i, err)
		}
	}
	if len(c.founds) > len(expected) {
		t.Errorf("%s: much founds: expected %d actually %d: next is %s",
			name, len(expected), len(c.founds), c.founds[len(expected)])
	}
}

const dict0 = `
ユーザー: [ ユーザ ]

サーバー:
  - サーバ

何:
  - なに

どんなに:

foo bar:
  - foobar
`

func parseDict(t *testing.T, s string) Dict {
	d, err := readDict(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestFind(t *testing.T) {
	d := parseDict(t, dict0)

	testFind(t, d, "empty", nil, ``)

	testFind(t, d, "basic", []found{
		{6, 18, "ユーザー", ""},
		{39, 48, "ユーザ", "ユーザー"},
		{69, 81, "サーバー", ""},
		{102, 111, "サーバ", "サーバー"},
		{141, 151, "ユーザ", "ユーザー"},
		{206, 219, "ユーザー", ""},
		{296, 310, "ユーザ", "ユーザー"},
		{387, 404, "ユーザー", ""},
		{459, 468, "サーバ", "サーバー"},
	}, `このユーザーはOKです。
このユーザはNGです。
このサーバーはOKです。
このサーバはNGです。
行マタギのユー
ザを検出できるかのチェック。
行マタギのユー
ザーを検出できるかのチェック。
    インデントと行マタギのユー
    ザを検出できるかのチェック。
    インデントと行マタギのユー
    ザーを検出できるかのチェック。
最終行のNG:サーバ検出`)

	testFind(t, d, "first line + LF", []found{
		{0, 10, "ユーザ", "ユーザー"},
	}, `ユー
ザ(先頭の行マタギNG検出)`)

	testFind(t, d, "many LF", []found{
		{0, 13, "ユーザー", ""},
		{114, 123, "ユーザ", "ユーザー"},
	}, `ユーザ
ー(改行だけの場合は、連続しているとみなされて検出されるべきではない)


ユーザ

ー(空行を挟んだ場合は、分離しているとみなされて検出されるべき)`)

	testFind(t, d, "NANI workaround", []found{
		{24, 30, "なに", "何"},
		{68, 80, "どんなに", ""},
	}, `# 検出されるべき
なに

# 検出されるべきではない
どんなに`)

	testFind(t, d, "spaces in words", []found{
		{24, 30, "foobar", "foo bar"},
		{68, 75, "foo bar", ""},
		{113, 120, "foo bar", ""},
	}, `# 検出されるべき
foobar

# 検出されるべきではない
foo bar

# 検出されるべきではない
foo
bar`)
}

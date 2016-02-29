package main

import "fmt"

// Found represents a found word.
type Found struct {
	Begin int
	End   int
	Word  *Word
}

func (f *Found) String() string {
	if f.Word.Fix != nil {
		return fmt.Sprintf("Found{Begin:%d, End:%d, Text:%q, Fix:%q}",
			f.Begin, f.End, f.Word.Text, *f.Word.Fix)
	}
	return fmt.Sprintf("Found{Begin:%d, End:%d, Text:%q}",
		f.Begin, f.End, f.Word.Text)
}

// OK returns true when there no fix (it is a correct word).
func (f *Found) OK() bool {
	return f.Word.Fix == nil
}

func (f *Found) IsBegin(offset int) bool {
	return f != nil && offset == f.Begin
}

func (f *Found) In(offset int) bool {
	return f != nil && offset >= f.Begin && offset < f.End
}

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

// IsBeginAndFix returns true if found's begin matches with offset and have
// valid fix.
func (f *Found) IsBeginAndFix(offset int) bool {
	return f != nil && offset == f.Begin && f.Word.Fix != nil
}

// In returns true if offset is between begin and end.
func (f *Found) In(offset int) bool {
	return f != nil && offset >= f.Begin && offset < f.End
}

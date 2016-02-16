package main

import "fmt"

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

func (f *Found) OK() bool {
	return f.Word.Fix == nil
}

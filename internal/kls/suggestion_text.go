package kls

import "fmt"

func (kls *KLS) printSuggestion(format string, message ...interface{}) {
	kls.suggestionText.Clear()
	fmt.Fprintf(kls.suggestionText, format, message...)
	kls.app.Draw()
}

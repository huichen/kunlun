package kls

import "fmt"

// 工具函数，打印提示字符串（搜索框右侧的文案）
func (kls *KLS) printSuggestion(format string, message ...interface{}) {
	kls.suggestionText.Clear()
	fmt.Fprintf(kls.suggestionText, format, message...)
	kls.app.Draw()
}

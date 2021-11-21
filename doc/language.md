多语言支持
======

#### 语言检测

我们使用 [go-enry](https://github.com/go-enry/go-enry) 做语言检测，这个库支持数百种语言和文件类型的识别，速度和准确度都比较高。

如果你只想索引某些或者不想索引某些语言，在代码中你可以通过 [IndexWalkerOptions.AllowedCodeLanguages 和 DisallowedCodeLanguages](https://github.com/huichen/kunlun/blob/master/pkg/types/walker_options.go) 选项进行控制。


#### 语言搜索过滤

通过表达式修饰词，比如

```
lang:java
lang:java or lang:cpp
-lang:python
```

见 [搜索语言](/doc/language.md)
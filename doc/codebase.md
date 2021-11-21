代码目录结构
======

随着新代码的提交，下面的目录结构可能不是最新的，仅供参考。

```
api
    - rest：        KWS 服务接口

cmd
    - benchmark：   性能评测程序
    - examples：    样例程序
    - kls：         KLS 命令行瑞士军刀程序
    - kws：         KWS 搜索网站服务

doc                 文档

internal            所有不对外暴露的内部库放这里
    - api           KWS 接口实现
    - ctags         ctags 解析器
    - indexer       索引器代码
    - kls           KLS 实现代码
    - ngram_index   三元组索引核心代码
    - query         搜索表达式解析器
    - ranker        排序器代码
    - resource      KWS 使用的资源句柄
    - searcher      搜索器代码
    - util          工具函数
    - walker        遍历器代码（索引前遍历仓库和文件）

pkg                 所有对外暴露的 SDK 放这里
    - engine        引擎
    - log           日志器
    - types         内外部公用的类型放这里

vendor              三方库
```

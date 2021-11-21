使用 ctags 支持符号检索
======

如果你希望在搜索表达式中支持 [sym: 类型的查询](/doc/query.md)，那么需要编译安装 [universal-ctags](https://github.com/universal-ctags/ctags)

编译时需要打开两个选项，以支持 json 输出和 [seccomp](https://en.wikipedia.org/wiki/Seccomp) 安全特性（如果希望在 Linux 上支持 ctags 那么 seccomp 是必须选项）

```
./configure --enable-json --enable-seccomp
```
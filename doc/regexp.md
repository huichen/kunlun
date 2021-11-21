正则表达式查询
======

正则表达式查询在所有查询模式中是时间复杂度最高的，比简单字符串查询的时间高一个数量级左右，我们对正则表达式查询做了深度优化，使得绝大多数正则表达式的查询时间降低到和简单字符串相同的程度。下面是几个优化点：

#### 单串提取预查寻

比如下面的正则表达式

```
cpu.*cache
```

我们的做法是将表达式中的简单字串比如 cpu 和 cache 提取出来，然后对 cpu 和 cache 分别查询，再将得到的文档 ID 做交集合并，从而极大减少待匹配候选集。

#### 单串查询优化

在查询单串 cache 时，因为这个词可以拆分成不同三元组，比如

```
cac
ach
che
```

我们在所有候选三元组中根据在文档中出现的频率，找到出现频率最低的两个三元组，然后按照取交集再做局部匹配的方式，从而大大加快了长单串的查询速度。

#### 行内表达式查询

多数表达式只搜索一行之内的内容，在上面“单串提取预查寻”中已经得到了匹配单串出现的行位置，可以利用这个信息进一步降低做正则表达式匹配的文本大小，因为正则表达式匹配的时间复杂度通常和待匹配字串的长度成正比，这个优化可以显著加快速度。

#### 表达式联合优化

在查询单串、正则混合类型的搜索表达式时，通常可以利用单串结果做优化，比如下面的表达式

```
cpu.*cache AND hit
```

我们可以先查询 hit 满足的文档，然后利用这个文档范围缩小 cpu.*cache 的匹配范围。

#### 并发查询

每次正则表达式我们都启动了多个线程做并发（goroutines 干这个太擅长了），这可以进一步降低总延迟。
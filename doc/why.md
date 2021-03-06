为什么要有昆仑引擎
======


#### 代码搜索的重要性

代码是软件公司最重要的无形资产，软件工程师每天都在企业内的代码仓库上花费大量时间做阅读、创作和维护工作。

代码搜索是工程师日常的重要工作之一，根据这篇 [Google 的论文](https://static.googleusercontent.com/media/research.google.com/zh-CN//pubs/archive/43835.pdf)统计，工程师平均每天要做 5.3 次代码搜索，而其中 34% 的搜索是为了寻找可供参考的代码样例。可以说代码搜索功能是提高工程师生产力的重要工具。

对单代码仓库，IDE 中自带的功能通常可以解决代码搜索问题，但对大规模团队协作的多仓库，或者复杂搜索条件，IDE 无法满足需求，需要独立的搜索引擎支持。

#### 昆仑引擎的目标

我们希望通过昆仑项目，**给企业提供数亿行代码级、支持复杂查询表达式、百毫秒内响应延迟的实时代码搜索能力**。

目前开源的代码搜索（比如 zoekt、opengrok、livegrep 等）均不能很好满足这个需求，有的对查询语言支持较弱（比如 opengrok），有的对复杂正则表达式查询速度很慢（比如 zoekt），有的无法索引这么大规模的数据。

因此我们参考了[三元组索引架构](https://swtch.com/~rsc/regexp/regexp4.html)，并做了大量工程优化，推出了昆仑搜索引擎。

我们希望通过这个项目，给企业级的代码搜索打下一个坚实的基础，希望通过一个好用的代码搜索工具，切实地提高软件工程师的工作效率。

昆仑欢迎商业使用，有任何问题或者合作意愿可以[联系我们](/doc/feedback.md)。
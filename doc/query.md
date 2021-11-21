### 昆仑支持的搜索语言

**cpu**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含 "cpu" 的文档

**"CPU cache"**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含短语 "CPU cache" 的文档

**cpu.\*name**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜匹配正则表达式 "cpu.*name" 的文档

**"cpu\d{3} name"**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜匹配正则表达式 "cpu\d{3} name" 的文档

**cpu cache**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含 "cpu" 并且包含 "cache" 的文档

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;同 cpu AND cache

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;同 cpu and cache

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;同 (cpu and cache)

**cpu or cache**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含 "cpu" 或者 "cache" 的文档

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;同 cpu OR cache

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;同 (cpu or cache)

**cpu -cache**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含 "cpu" 并且不包含 "cache" 的文档

**cpu -(cache or miss)**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含 "cpu" 并且不包含 "cache" 或 "miss" 的文档

**cpu -"cache miss"**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜索包含 "cpu" 但不含 "cache miss" 短语的文档

**cpu cache or hit miss**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;等价于搜 (cpu AND cache) OR (hit AND miss)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- 优先级高于 AND 高于 OR，可以用这三个操作符加括号组合任意深度的表达式

**CPU cache.\*name case:yes**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含 "CPU" 并且匹配 "cache.*name" 的文档，两者都大小写敏感。默认不区分大小写

**cpu file:admin**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含 "cpu" 并且文件名中包含 "admin" 的文档

**cpu (file:api.\*doc or file:admin)**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含 "cpu" 并且文件名匹配正则表达式 "api.*doc" 或者包含 "admin" 的文档

**cpu lang:java**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含 "cpu" 的 java 代码

**cpu -lang:java**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含 "cpu" 的非 java 代码

**cpu -(lang:java or lang:python)**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含 "cpu" 的 java/python 之外的代码

**file:\.cpp$**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜索文件名以 ".cpp" 结尾的文档

**cpu -file:admin.\*java -file:web**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜包含 "cpu" 但文件名不匹配正则表达式 "admin.*java" 也不包含 "web" 的文档

**sym:data**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜索包含符号 "data" 的文档，sym:不可以作用在正则表达式上

**cpu.\*name repo:web.\*service**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜索匹配正则表达式 "cpu.*name"，同时所在仓库名匹配正则表达式 "web.*service" 的文档

**repo:web.\*service**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜索仓库名匹配正则表达式 "web.*service" 的代码仓库

**repo:web.\*service -repo:admin**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;搜索仓库名匹配正则表达式 "web.*service" 但名称中不包含 "admin" 的代码仓库

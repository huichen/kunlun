访问权限控制
======

昆仑支持“人+仓库组合”粒度的代码控制，可以给不同的访问者返回只有该用户能访问到的仓库里的搜索结果，同时保持检索速度不受影响。

具体做法是通过 SearchRequest.ShouldRecallRepo 传入一个钩子函数，这个函数接收仓库 ID 作为参数，请在外部根据访问者身份生成这个钩子，搜索器会实时调用该函数确定是否在搜索结果中返回；外部预先生成的 repoID 请通过 [EngineOptions.RepoRemoteURLToIDMap](https://github.com/huichen/kunlun/blob/master/pkg/types/engine_options.go) 在引擎启动时传入。

**注意**：请使用本地缓存优化这个函数的速度，因为最差的情况下，检索时会对所有仓库调用这个函数，最好的方案是在这个函数中保存一个所有仓库 ID 的内存中的白名单，这个白名单在引擎查询时不能有写操作，否则会有读写冲突。
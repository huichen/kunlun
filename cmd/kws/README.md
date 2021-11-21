昆仑代码网页搜索（Kunlun Web Search）
=======

KWS 提供了一个代码搜索的 Web 服务，方便你和你的团队通过网页界面做代码查询和浏览。

![](https://github.com/huichen/kunlun/blob/master/doc/kws.png)

#### 如何启动 KWS

1、编译前端页面：进入 cmd/kws/fe 目录，输入

```
npm install
npm run build
```

构建好的页面会放在 cmd/kws/fe/dist 目录下

2、启动服务：进入 cmd/kws 目录，输入

```
go run main.go -repo_folders <你的代码目录列表>
```

其中代码目录列表用半角逗号分隔，中间无空格。

#### KWS 参数配置

KWS 支持一系列的启动参数，可以通过下面的命令行得到

```
go run main.go -help
```

默认参数通常足够好。

#### 关于代码仓库

目前只支持 git 仓库，建议把所有需要检索的仓库 git clone 到一个目录下，然后 -repo_folders 指向这个目录即可。
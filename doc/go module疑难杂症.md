在Go 1.13中，我们可以通过GOPROXY来控制代理，以及通过GOPRIVATE控制私有库不走代理。

设置GOPROXY代理：

```shell
go env -w GOPROXY=https://goproxy.cn,direct
```

设置GOPRIVATE来跳过私有库，比如常用的Gitlab或Gitee，中间使用逗号分隔：

```sh
go env -w GOPRIVATE=*.gitlab.com,*.gitee.com
```

//go 设置代理

```
go env -w GOPROXY=https://goproxy.cn,direct  //设置代理

go env -w GOPRIVATE=*.gitlab.com,*.gitee.com  //
```



Go 1.13提供了GOSUMDB环境变量用于配置Go校验和数据库的服务地址（和公钥），其默认值为”sum.golang.org”，这也是Go官方提供的校验和数据库服务(大陆gopher可以使用sum.golang.google.cn)。

出于安全考虑，建议保持GOSUMDB开启。但如果因为某些因素，无法访问GOSUMDB（甚至是sum.golang.google.cn）；这时候可以用下面的命令将关闭掉sum校验功能。

```shell
 //关闭掉gosumdb的远程校验，仅能使用本地的go.sum进行包的校验和校验了。

go env -w GOSUMDB=off    

 //国内的sum验证服务

go env -w GOSUMDB="sum.golang.google.cn"  
```



当我们需要访问一些内网的私有的仓库地址的时候，这个时候是不需要使用代理的，这时候，设置GONOPROXY来忽略代理。

```shell
go env -w GONOPROXY=*.gitlab.com
```



背景介绍:

在内网部署了一个 gitlab 的服务，域名是 [http://git.abc.com](http://git.abc.com/) 。他们的程序员需要使用 Go Module 的功能。 由于有些公共代码是公司内部的，而有些依赖则是 Github 上开源的，另外为了提高 Go Module 下载速度，他们还打算用xxx云的加速服务。

问题来了，Go Module 默认支持的是 https 的域名，所以内网的 gitlab 服务就比较麻烦，也就是说内网的 gitlab 的库下载的时候要满足两点需求，一个是支持 http 协议，另外一个是不使用xxx的代理。

(

案例:例如因为上面这种情况，在使用过程中，出现过的问题：

reading https//sum.golang.google.cn/lookup/gitlab.abc.com/xxx：410 Gone

server response : not found   gitlab.abc.com/xxx：unrecognized import path https fetch: Get "https://gitlab.abc.com/xxx?go-get=1": EOF

)

方案:

我们需要引入两个环境变量来支持私有库的 http 协议和忽略使用xxx云的代理

```shell
go env -w GONOPROXY=git.abc.com

go env -w GOINSECURE=git.abc.com   //设置使用http，而不是https访问
```



其中允许非安全下载，主要是针对没有HTTPS的HTTP路径：

go get -insecure   
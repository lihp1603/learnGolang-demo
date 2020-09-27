## gRPC开发环境搭建

### 1. protobuf相关环境

golang 使用protobuf其实网上已经有很多的讲解了，
这篇主要讲我自己开始搞的时候有点迷糊的地方，可能会帮到其他跟我一样新手的人。

首先下载protobuf的 编译器protoc 地址：记得下载对版本
window:
下载 protoc-3.6.1-win32.zip
解压后，把bin目录下的protoc.exe文件复制到到你的%GOPATH%目录下

接着我们要获取到protobuf的go插件

```tex
go get -u github.com/golang/protobuf/protoc-gen-go
```

如果获取不到的话 可以自己去gitbub上面 clone或是下载
记得下载后要复制到github.com/golang这个目录下，然后进行手动安装:

```tex
go install github.com\golang\protobuf\protoc-gen-go
```

这样就会在%GOPATH%目录下再生成一个protoc-gen-go.exe文件。

(有了protoc.exe和protoc-gen-go.exe这两个，我们就能处理proto文件了，让我们自己定义的.proto文件生成对应想要的程序代码文件。)

如果你直接把proto文件跟protoc.exe 放在同一个目录 那么进入cmd 执行protoc  *.proto  --go_out=.  就可以了。
另外可以指定proto目录和生成目录,执行:

```tex
protoc --proto_path=你proto文件的目录 --go_out=生成go文件的目录 你的proto文件.proto
```

(原文链接：https://blog.csdn.net/hawu_hao/article/details/82909951)

### 2. gRPC环境

我这里使用go 1.14版本，同时也尝试使用go module功能。

- 2.1  gRPC安装

```go
go get -u google.golang.org/grpc
```

针对gRPC安装失败，官方上面的方法和建议:

The `golang.org` domain may be blocked from some countries. `go get` usually produces an error like the following when this happens:

```
$ go get -u google.golang.org/grpc
package google.golang.org/grpc: unrecognized import path "google.golang.org/grpc" (https fetch: Get https://google.golang.org/grpc?go-get=1: dial tcp 216.239.37.1:443: i/o timeout)
```

To build Go code, there are several options:

- Set up a VPN and access google.golang.org through that.

- Without Go module support: `git clone` the repo manually:

  ```
  git clone https://github.com/grpc/grpc-go.git $GOPATH/src/google.golang.org/grpc
  ```

  You will need to do the same for all of grpc's dependencies in `golang.org`, e.g. `golang.org/x/net`.

- With Go module support: it is possible to use the `replace` feature of `go mod` to create aliases for golang.org packages. In your project's directory:

  ```
  go mod edit -replace=google.golang.org/grpc=github.com/grpc/grpc-go@latest 
  go mod tidy   
  go mod vendor
  go build -mod=vendor
  ```

这里我使用git clone的方法，然后将grpc放到github.com/grpc/grpc-go下,然后使用go mod的方法就OK了。



- 2.2  简单介绍一下go mod方法：

```tex
go mod init 
生成go.mod文件
go mod edit -replace=google.golang.org/grpc=github.com/grpc/grpc-go@latest  
编辑和替换go.mod文件
go mod tidy      
会检测该文件夹目录下所有引入的依赖,写入 go.mod 文件
go mod vendor    
执行此命令,会将刚才下载至 GOPATH 下的依赖转移至该项目根目录下的 vendor(自动新建) 文件夹下
go build -mod=vendor
编译项目
```

这里说一下使用代理的方法:

```tex
 即在执行前先设置好环境变量“GOPROXY”和“GO111MODULE”： 

 export GO111MODULE=on 

 export GOPROXY=https://goproxy.io或者export GOPROXY=https://goproxy.cn 

 对于1.13及以上版本，可直接如下这样： 

 go env -w GOPROXY=https://goproxy.cn,direct 

 推荐使用七牛云的“goproxy.cn”，因为“goproxy.io”也不一定可用
```

note:  这里说一下go mod下，引用本地package的存在的问题。

(可参考:https://www.liwenzhou.com/posts/Go/import_local_package_in_go_module/)

文章中介绍的方法，基本可用，我们使用go mod init以后，会生成一个go.mod文件，然后接着使用go mod tidy去更新和下载依赖包的时候，如果是本地包的话，可能会存在问题：

例如，我这里使用 import github.com/lihp1603/gRPC_demo/xlive_proto 这个package，但是这个package是我自己本地写的，而且还没有上传的情况下，我们可以手动打开go.mod文件，然后手动修改,添加内容:

require github.com/lihp1603/gRPC_demo/xlive_proto v0.0.0

replace github.com/lihp1603/gRPC_demo/xlive_proto => ../xlive_proto  

这样就将这个package改为了本地的引用。

这里需要注意，使用go mod功能以后，并去搜索go path下的package。












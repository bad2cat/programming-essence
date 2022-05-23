#### Golang 交叉编译

GO 支持交叉编译，也就是说，可以用一个平台下的 `Golang `编译器，生成另外一个平台下的可执行程序

平时一般是在 `Linux`下编程，并且大部分的时间也是使用`Docker`镜像，所以一般不需要不会在意这个交叉编译的功能

但是，这两天公司另外一个部门突然让生成一个`exe`文件，刚开始想的是找个安装`go`语言的`windows`平台，然后将`git`仓库上的代码拉取下来执行`go build`就好了；后来发现，使用交叉编译实在是太方便了，只需要执行一行编译的命令就搞定了

在`linux`下生成`Mac`和`Window`的可执行程序

```
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
```

在`window`下生成`Mac`和`Linux`的可执行程序

```
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build
```

在`Mac`下生成`Linux`和`Windows`下的可执行程序

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
```

心想这样太简单了，这部分分中的事么，结果真是啪啪打脸啊

因为，在程序中用到了`sqlite`，而`sqlite`中又需要开启`CGO`，也就是说，需要在`Go`代码中用到`c`语言编程，或者是调用`C`代码的函数库

而`golang`的交叉编译只能编译不同平台的`go`代码，但是编译不了不同平台的`C`代码，所以还需要安装一个指定平台的`C`代码的编译器

因为，打算是在`Linux`下编译生成`windows`下`64`的可执行程序，所以，就在`Linux`下安装了一个`windows`下的`C`编译器，命令如下：

```
sudo apt install mingw-w64
```

这条命令安装了如下的两个编译器：

```
32位windows编译器i686-w64-mingw32-gcc和i686-w64-mingw32-g++
64位windows编译器x86_64-w64-mingw32-gcc和x86_64-w64-mingw32-g++
```

安装之后，就可以指定用这个编译器编译了，命令如下：

```
export GOOS=windows
export GOARCH=amd64
export CC=x86_64-w64-mingw32-gcc  //此处是指定 64 位的编译器，你可以指定 32 位的
export CXX=x86_64-w64-mingw32-g++
export CGO_ENABLED=1
```

配置`env`之后，再执行`go build`就`ok`了；当然，你也可以那这写成一个脚本文件：

```
#!/bin/bash
export GOOS=windows
export GOARCH=amd64
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++
export CGO_ENABLED=1
go build .
```

之后需要交叉编译时，直接执行这个脚本就好了

注：因为`golang`默认是开启`cgo`的，你可以用`go env`看到`CGO_ENABLED="1"`，所以，如果你程序中没用到`c`编程或者调用`C`的函数库，那最后在编译的时候指定`CGO_ENABLED=0`，这样交叉编译就不需要这么复杂了

参考地址：https://blog.furry.top/archives/5/
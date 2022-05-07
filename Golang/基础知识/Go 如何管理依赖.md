### Golang 依赖管理

在 `go`中有两个非常重要的环境变量：

- `GOROOT`： `golang`的安装路径，Linux 下默认会安装在`/usr/local/go`之下
- `GOPATH`：存放`SDK`以外的第三方类库；收藏的可复用的代码，包含三个子目录： -- src : 存放项目源码文件 -- pkg : 编译后的包文件 -- bin ：编译后生成的可执行文件

#### Go Module

从 `Go1.11`版本开始，引入了`Go Module`功能，在程序中只有一个`go.mod`文件，存放依赖列表，而依赖的具体的包默认会下载到`GOPATH/pkg/mod`目录下面

在运行 `go build`时，优先引用`go.mod`文件中的依赖，所以`Vendor`目录逐渐的被替代了，消失在了大部分的项目当中

##### `go mod`的使用

在`go 1.13`以后，`go mod`就是默认的依赖管理工具，在项目中使用如下的方式创建`go.mod`文件夹管理依赖

```
go mod int
```

这样整个项目就可以使用`go mod`来进行管理了

下面是 `go mod `的一些常用命令：

```
go mod download 下载模块到本地缓存，缓存路径是 $GOPATH/pkg/mod/cache
go mod edit 是提供了命令版编辑 go.mod 的功能，例如 go mod edit -fmt go.mod 会格式化 go.mod
go mod graph 把模块之间的依赖图显示出来
go mod init 初始化模块（例如把原本dep管理的依赖关系转换过来）
go mod tidy 增加缺失的包，移除没用的包
go mod vendor 把依赖拷贝到 vendor/ 目录下
go mod verify 确认依赖关系
go mod why 解释为什么需要包和模块
```

`go mod`中存在的一些问题：

1. 依赖包的地址失效或者下载失败

   ```
   此时可以使用 go.mod 文件中的 replace 来替换这个包
   ```

2. 之前的项目如何使用`go mod`管理

   ```
   首先把项目移出$GOPATH/src目录，然后运行 go mod init module_name，最后执行 go build 即可
   ```

#### Go Vendor

`Vendor` 目录是`Golang1.5`版本引入的，为项目提供了一种离线保存第三方依赖的方法；使用`vendor`的项目，搜索依赖的顺序如下：

```
当前包下的 vendor 目录
向上级目录查找，直到找到 src 下的 vendor 目录
在 GOROOT 目录下查找
在 GOPATH 下面查找依赖包
```

###### 项目中使用`go vendor`

安装

```
go get -u -v github.com/kardianos/govendor
```

初始化

```
cd xxx
govendor init
```

初始化完成后，项目目录中会生成一个vendor文件夹，包含一个vendor.json文件，json文件中包含了项目所依赖的所有包信息

```text
{
 "comment": "",
 "ignore": "test",
 "package": [],
 "rootPath": "govendor-test"
}
```

###### `govendor`常用命令

- 将已被引用且在 $GOPATH 下的所有包复制到 vendor 目录

  ```
  govendor add +external
  ```

- 仅从 $GOPATH 中复制指定包

  ```
  govendor add gopkg.in/yaml.v2
  ```

- 列出一个包被哪些包引用

  ```
  govendor list -v fmt
  ```

- 从远程仓库添加或更新某个包(**不会**在 `$GOPATH` 也存一份)

  ```
  govendor fetch golang.org/x/net/context
  ```

- 安装指定版本的包

  ```
  govendor fetch golang.org/x/net/context@a4bbce9fcae005b22ae5443f6af064d80a6f5a55
  govendor fetch golang.org/x/net/context@v1   # Get latest v1.*.* tag or branch.
  govendor fetch golang.org/x/net/context@=v1  # Get the tag or branch named "v1".
  ```

- 只格式化项目自身代码(`vendor` 目录下的不变动)

  ```
  govendor fmt +local
  ```

- 拉取所有依赖的包到 `vendor` 目录(包括 `$GOPATH` 存在或不存在的包)

  ```
  govendor fetch +out
  ```

- 包已在 `vendor` 目录，但想从 `$GOPATH` 更新

  ```
  govendor update +vendor
  ```

各子命令详细用法可通过 `govendor COMMAND -h` 或阅读 `github.com/kardianos/govendor/context` 查看源码包如何实现的。

| 子命令  | 功能                                                         |
| :-----: | ------------------------------------------------------------ |
|  init   | 创建 `vendor` 目录和 `vendor.json` 文件                      |
|  list   | 列出&过滤依赖包及其状态                                      |
|   add   | 从 `$GOPATH` 复制包到项目 `vendor` 目录                      |
| update  | 从 `$GOPATH` 更新依赖包到项目 `vendor` 目录                  |
| remove  | 从 `vendor` 目录移除依赖的包                                 |
| status  | 列出所有缺失、过期和修改过的包                               |
|  fetch  | 从远程仓库添加或更新包到项目 `vendor` 目录(不会存储到 `$GOPATH`) |
|  sync   | 根据 `vendor.json` 拉取相匹配的包到 `vendor` 目录            |
| migrate | 从其他基于 `vendor` 实现的包管理工具中一键迁移               |
|   get   | 与 `go get` 类似，将包下载到 `$GOPATH`，再将依赖包复制到 `vendor` 目录 |
| license | 列出所有依赖包的 LICENSE                                     |
|  shell  | 可一次性运行多个 `govendor` 命令                             |

###### govendor 状态参数

|   状态    | 缩写 | 含义                                                 |
| :-------: | :--: | ---------------------------------------------------- |
|  +local   |  l   | 本地包，即项目内部编写的包                           |
| +external |  e   | 外部包，即在 `GOPATH` 中、却不在项目 `vendor` 目录   |
|  +vendor  |  v   | 已在 `vendor` 目录下的包                             |
|   +std    |  s   | 标准库里的包                                         |
| +excluded |  x   | 明确被排除的外部包                                   |
|  +unused  |  u   | 未使用的包，即在 `vendor` 目录下，但项目中并未引用到 |
| +missing  |  m   | 被引用了但却找不到的包                               |
| +program  |  p   | 主程序包，即可被编译为执行文件的包                   |
| +outside  |      | 相当于状态为 `+external +missing`                    |
|   +all    |      | 所有包                                               |

支持状态参数的子命令有：`list`、`add`、`update`、`remove`、`fetch`

###### `govendor`的缺点

- 使用`vendor`的项目需要把该项目所有的依赖都下载到这个项目中的`vendor`目录下，这就会导致项目变的很庞大

- `vendor`不区分依赖包的版本，这就意味着可能出现不同的环境依赖了不同版本的包，很有可能导致编译出错

#### Go mod vendor

为了能够在离线的环境中也能够部署和应用程序，就保存了`vendor`目录；然后，执行如下的命令：

```
go mod vendor
```

就可以将当前程序的依赖拷贝到`vendor`目录下，当程序中的包使用`go mod`下载不下来的时候，就可以引用`vendor`目录下的依赖了
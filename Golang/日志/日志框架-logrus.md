#### Gin 使用 logrus 打印日志

##### logrus 是什么

[logrus](https://github.com/sirupsen/logrus)是一个比较常用的日志框架，它的 API 和 `Go`标准库中的`log`是完全兼容的，所以可用`logrus`替换`go`标准库中的`log`；`logrus`目前处于维护状态，不会在添加新的特性了

##### 如何使用`logrus`

下面是使用`logrus`打印日志的一个案例：

```go
package main

import (
  "os"
  log "github.com/sirupsen/logrus"
)

func init() {
  // Log as JSON instead of the default ASCII formatter.
  log.SetFormatter(&log.JSONFormatter{})

  // Output to stdout instead of the default stderr
  // Can be any io.Writer, see below for File example
  log.SetOutput(os.Stdout)

  // Only log the warning severity or above.
  log.SetLevel(log.WarnLevel)
}

func main() {
  log.WithFields(log.Fields{
    "animal": "walrus",
    "size":   10,
  }).Info("A group of walrus emerges from the ocean")

  // A common pattern is to re-use fields between logging statements by re-using
  // the logrus.Entry returned from WithFields()
  contextLogger := log.WithFields(log.Fields{
    "common": "this is a common field",
    "other": "I also should be logged always",
  })

  contextLogger.Info("I'll be logged with common and other field")
  contextLogger.Info("Me too")
}
```

上面的`log.SetFormatter`就是设置日志的输出格式，`logrus`有两种输出格式：

```go
log.JSONFormatter()  // 以 json 格式输出
log.TextFormatter()  // 以文本的方式输出
```

`log.SetOutput(os.Stdout)`是设置日志输出的位置，默认是`os.Stderr`

`log.SetLevel(log.WarnLevel)`设置日志的级别，`logrus`支持七个级别的日志，分别为：`TraceLevel,DebugLevel,InfoLevel,WarnLevel,ErrorLevel,FatalLevel,PanicLevel`，这是日志级别从低到高的排序，设置低级别的日志，比如：`DebugLevel`，那么会输出该级别以上级别的日志，但是不会输出该级别以下的日志，因此`DebugLevel`就不会输出`TraceLevel`的日志了

##### Gin 中使用 logrus 输出日志

`Gin`是通过中间件打印日志的，在`Gin`中通过如下的方式添加中间件，就能输出`Gin`的日志了

```go
func main() {
    // Disable Console Color, you don't need console color when writing the logs to file.
    gin.DisableConsoleColor()

    // Logging to a file.
    f, _ := os.Create("gin.log")
    gin.DefaultWriter = io.MultiWriter(f)

    // Use the following code if you need to write the logs to file and console at the same time.
    // gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

    router := gin.Default()
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    router.Run(":8080")
}
```

`gin.DisableConsoleColor`就是关闭输出日志时带的颜色功能

`router := gin.Default()`中包含了两个中间件，如下：

```
func Default() *Engine {
	debugPrintWARNINGDefault()
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}
```

一个是`Logger`中间件就是打印`Gin`的日志的组件；一个是`Recovery`组件，防止程序崩溃的组件

`Gin`打印的日志格式是文本格式，如何你想打印`json`格式的日志，可以使用`logrus`；首先关闭`Gin`的日志，所以就不能通过`gin.Default()`获取`engine`了，代码如下

```go
func main() {
    router := gin.New()
    router.Use(gin.Recovery())
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    router.Run(":8080")
}
```

下面自己实现一个日志的中间件，然后在`gin`中使用这个中间件就好了，日志中间件的代码如下：

```go
func LogMiddleWare() gin.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	logger.Out = os.Stdout

	return func(c *gin.Context) {
		c.Next()
		method := c.Request.Method
		reqUrl := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		logger.WithFields(logrus.Fields{
			"method":      method,
			"uri":         reqUrl,
			"status_code": statusCode,
			"client_ip":   clientIP,
		}).Info()
	}
}
```

接下来就是在`gin`中使用这个中间件：

```go
func main() {
    router := gin.New()
    router.Use(gin.Recovery())
    router.User(LogMiddleWare())
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    router.Run(":8080")
}
```

接下来看下日志打印的结果：

```json
{"client_ip":"0.0.0.0","level":"info","method":"GET","msg":"","status_code":200,"time":"2022-04-26T18:19:20+08:00","uri":"/ping"}
```

可以看到打印出来的日志已经变成了`json`格式了

##### End

当然，`gin`的`log`组件也可以自定义打印出来的字段，比如：

```go
func main() {
	router := gin.New()

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.Run(":8080")
}
```

但这种输出的格式也不是`json`的，所以选择那种还是要根据自己的需求而定

参考链接：https://github.com/gin-gonic/gin#how-to-write-log-file
### HTTP Server

HTTP Server 就是用来处理客户端请求，然会返回处理结果的服务器，这是典型的 CS 结构

首先，HTTP Server 在某个站点监听请求，当客户端发送请求之后，服务端就会按照一定的规则就处理，处理完成返回结果

下面是一个简单的 Go HTTP Server:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

type tHandler struct {
	content string
}

func (t *tHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w,t.content)
}

func main() {
    http.HandleFunc("/",handler)
	http.Handle("/t",&tHandler{content:"Test"})
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}
```

使用 go 创建 http Server 比较简单，只有几行代码。上面虽然有两种注册的方式，但是其实在低层都是一种方式，下面先来看看 `http.Handle("/t",&tHandler{content:"Test"})`的源码

```go
func Handle(pattern string, handler Handler) { 
    DefaultServeMux.Handle(pattern, handler) 
}
```

这个 Handle 方法中调用了`DefaultServeMux.Handle(pattern,handler)`，这个`DefaultServeMux`是`ServeMux`类型的一个默认值，`ServeMux`的结构如下：

```go
type ServeMux struct {
	mu    sync.RWMutex
    //这个 map 中 key 就是 URL path 而 muxEntry 就是一个 handler
	m     map[string]muxEntry 
	es    []muxEntry // slice of entries sorted from longest to shortest.
	hosts bool       // whether any patterns contain hostnames
}
```

`ServeMux`就是用来存放`path`对应的`handler`的。上面调用的`DefaultServeMux.Handle(pattern,handler)`就是用来注册路由的，以及这个路由对应的 `Handler`的

下面再来看看` http.HandleFunc("/",handler)`这个函数，下面是它的源码：

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern.
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerFunc(handler))
}
```

它首先是通过`DefaultServeMux`调用了`ServerMux`的`HandlerFunc`，而在这个方法中，同样也是调用了`ServeMux.Handle`方法。

也就是说，上面两种注册路由的方式其实都是相同的，它们都是通过`DefaultServeMux`来完成路由的注册功能的

下面注重来看一下`ServeMux`这部分的代码，首先看一下他们调用的`mux.Handle()`这个方法，下面是具体的代码：

```go
// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *ServeMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}
	// 如果是第一个路由，则创建一个空的 map
	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	e := muxEntry{h: handler, pattern: pattern}
	mux.m[pattern] = e
	if pattern[len(pattern)-1] == '/' {
		mux.es = appendSorted(mux.es, e)
	}
	if pattern[0] != '/' {
		mux.hosts = true
	}
}
```

这一段的代码就是注册路由和路由对应的`Handler`，注册的路由是直接写入到`mux.m`这个`map`当中的，这个`map`的结构体如下：

```go
m := map[string]muxEntry
//muxEntry
type muxEntry struct{
    pattern string
    handler Handler
}
```

刚才发现一个问题，这个`map[string]muxEntry`的`key`已经存储了`pattern`为什么还要再这个`muxEntry`中还要存储这个`pattern`而且值还是一样的？

后面发现了，在`ServeMux`中，还有一个`es []muxEntry`的结构体，这个结构体是通过`pattern`的长短对最后一个字符为`/`的`handler`进行排序用的，目前还只能理解到这了

下面先介绍一下，什么是`http.Handler`。在`net`中定义了一个接口，如下所示：

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

凡是实现了这个接口的类型都可以称之为是`Handler`，并且在`http`中又定义了一个类型：

```go
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```

这个类型实现了这个接口，也就是说， HandlerFunc 也是一个 `Handler`，这样的话，就可以对指定的函数返回匿名函数，然后这个匿名函数声明为`HandlerFunc`的样子就好了，比如：

```go
func Hang() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer,"test")
	}
}
```

通过这个匿名函数就可以直接返回一个 `Handler`了，当然这种也可以通过闭包向里面传递参数

当执行完`mux.Handle()`这个方法之后，整个路由和对应 `handler`的注册就全部完成了

接下来就是通过`http.ListenAndServe(":8000", nil)`来监听客户端的请求了

### Http.Server

上面的部分是注册路由以及对应的`Handler`部分，下面就是通过`http.Server`来监听用户的请求了

http.Server 的代码如下所示：

```go
log.Fatal(http.ListenAndServe(":8000", nil))

// 源码
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```

监听的时候需要传入一个`addr`，默认就是`:80`，以及一个`handler`，如果传入`nil`，则这个`handler`就是`http.DefaultServeMux`

下面是`server.ListenAndServe()`的源码：

```go
func (srv *Server) ListenAndServe() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}
```

可以看到，`http.Server`是通过`ne.Listen`方法获取到一个`tcp`的`Listener`，这是一个接口里面就三个方法，

```go
type Listener interface {
	// Accept waits for and returns the next connection to the listener.
	Accept() (Conn, error)

	// Close closes the listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Close() error

	// Addr returns the listener's network address.
	Addr() Addr
}
```

因为`go`支持通过不同的方式通信，上面的这个就是传入`tcp`，也就是去获取到了`tcp`的`Listener`，这部分代码是`net`中的，也就是更加底层的，暂时先不说

现在来看下`srv.Serve(ln)`这个方法，这个就是真正的接收`http`请求，然后调用对应的`Handler`进行处理的请求

```go
func (srv *Server) Serve(l net.Listener) error {
	if fn := testHookServerServe; fn != nil {
		fn(srv, l) // call hook with unwrapped listener
	}

    // 这部分就是通过调用 Listener.Close() 来关闭连接的，放到一个 onceCloseListen 中就是为了防止多个关闭
	origListener := l
	l = &onceCloseListener{Listener: l}
	defer l.Close()

    // 如果用到了 http2 这个就是用来初始化 http2 的配置的
	if err := srv.setupHTTP2_Serve(); err != nil {
		return err
	}
	
    // 添加当前的 Listener 到这个 Server.listerner 当中
	if !srv.trackListener(&l, true) {
		return ErrServerClosed
	}
	defer srv.trackListener(&l, false)

	baseCtx := context.Background()
	if srv.BaseContext != nil {
		baseCtx = srv.BaseContext(origListener)
		if baseCtx == nil {
			panic("BaseContext returned a nil context")
		}
	}

	var tempDelay time.Duration // how long to sleep on accept failure

	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
		rw, err := l.Accept()
		if err != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.logf("http: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		connCtx := ctx
		if cc := srv.ConnContext; cc != nil {
			connCtx = cc(connCtx, rw)
			if connCtx == nil {
				panic("ConnContext returned nil")
			}
		}
		tempDelay = 0
		c := srv.newConn(rw)
		c.setState(c.rwc, StateNew, runHooks) // before Serve can return
		go c.serve(connCtx)
	}
}
```

在`for`循环之前，首先是构造一个上下文`Context`，然后就是循环监听用户的请求，如下：

```go
rw, err := l.Accept()
```

当接收到客户端请求之后，就会构建一个 `http.conn`，它的结构如下：

```go
// A conn represents the server side of an HTTP connection.
type conn struct {
	// server is the server on which the connection arrived.
	// Immutable; never nil.
	server *Server

	// cancelCtx cancels the connection-level context.
	cancelCtx context.CancelFunc

	// rwc is the underlying network connection.
	// This is never wrapped by other types and is the value given out
	// to CloseNotifier callers. It is usually of type *net.TCPConn or
	// *tls.Conn.
	rwc net.Conn

	// remoteAddr is rwc.RemoteAddr().String(). It is not populated synchronously
	// inside the Listener's Accept goroutine, as some implementations block.
	// It is populated immediately inside the (*conn).serve goroutine.
	// This is the value of a Handler's (*Request).RemoteAddr.
	remoteAddr string

	// tlsState is the TLS connection state when using TLS.
	// nil means not TLS.
	tlsState *tls.ConnectionState

	// werr is set to the first write error to rwc.
	// It is set via checkConnErrorWriter{w}, where bufw writes.
	werr error

	// r is bufr's read source. It's a wrapper around rwc that provides
	// io.LimitedReader-style limiting (while reading request headers)
	// and functionality to support CloseNotifier. See *connReader docs.
	r *connReader

	// bufr reads from r.
	bufr *bufio.Reader

	// bufw writes to checkConnErrorWriter{c}, which populates werr on error.
	bufw *bufio.Writer

	// lastMethod is the method of the most recent request
	// on this connection, if any.
	lastMethod string

	curReq atomic.Value // of *response (which has a Request in it)

	curState struct{ atomic uint64 } // packed (unixtime<<8|uint8(ConnState))

	// mu guards hijackedv
	mu sync.Mutex

	// hijackedv is whether this connection has been hijacked
	// by a Handler with the Hijacker interface.
	// It is guarded by mu.
	hijackedv bool
}
```

这个`conn`就是服务端的一个`http`连接，构建完成之后就会启动一个`go routine`去处理这次请求；每个请求到来之后，服务端都会创建一个`go routine`去处理这个请求

下面就是通过`http.conn`处理请求的具体源码了：

```go
// Serve a new connection.
func (c *conn) serve(ctx context.Context) {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			c.server.logf("http: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
		}
		if !c.hijacked() {
			c.close()
			c.setState(c.rwc, StateClosed, runHooks)
		}
	}()

	if tlsConn, ok := c.rwc.(*tls.Conn); ok {
		if d := c.server.ReadTimeout; d > 0 {
			c.rwc.SetReadDeadline(time.Now().Add(d))
		}
		if d := c.server.WriteTimeout; d > 0 {
			c.rwc.SetWriteDeadline(time.Now().Add(d))
		}
		if err := tlsConn.HandshakeContext(ctx); err != nil {
			// If the handshake failed due to the client not speaking
			// TLS, assume they're speaking plaintext HTTP and write a
			// 400 response on the TLS conn's underlying net.Conn.
			if re, ok := err.(tls.RecordHeaderError); ok && re.Conn != nil && tlsRecordHeaderLooksLikeHTTP(re.RecordHeader) {
				io.WriteString(re.Conn, "HTTP/1.0 400 Bad Request\r\n\r\nClient sent an HTTP request to an HTTPS server.\n")
				re.Conn.Close()
				return
			}
			c.server.logf("http: TLS handshake error from %s: %v", c.rwc.RemoteAddr(), err)
			return
		}
		c.tlsState = new(tls.ConnectionState)
		*c.tlsState = tlsConn.ConnectionState()
		if proto := c.tlsState.NegotiatedProtocol; validNextProto(proto) {
			if fn := c.server.TLSNextProto[proto]; fn != nil {
				h := initALPNRequest{ctx, tlsConn, serverHandler{c.server}}
				// Mark freshly created HTTP/2 as active and prevent any server state hooks
				// from being run on these connections. This prevents closeIdleConns from
				// closing such connections. See issue https://golang.org/issue/39776.
				c.setState(c.rwc, StateActive, skipHooks)
				fn(c.server, tlsConn, h)
			}
			return
		}
	}

	// HTTP/1.x from here on.

	ctx, cancelCtx := context.WithCancel(ctx)
	c.cancelCtx = cancelCtx
	defer cancelCtx()

	c.r = &connReader{conn: c}
	c.bufr = newBufioReader(c.r)
	c.bufw = newBufioWriterSize(checkConnErrorWriter{c}, 4<<10)

	for {
		w, err := c.readRequest(ctx)
		if c.r.remain != c.server.initialReadLimitSize() {
			// If we read any bytes off the wire, we're active.
			c.setState(c.rwc, StateActive, runHooks)
		}
		if err != nil {
			const errorHeaders = "\r\nContent-Type: text/plain; charset=utf-8\r\nConnection: close\r\n\r\n"

			switch {
			case err == errTooLarge:
				// Their HTTP client may or may not be
				// able to read this if we're
				// responding to them and hanging up
				// while they're still writing their
				// request. Undefined behavior.
				const publicErr = "431 Request Header Fields Too Large"
				fmt.Fprintf(c.rwc, "HTTP/1.1 "+publicErr+errorHeaders+publicErr)
				c.closeWriteAndWait()
				return

			case isUnsupportedTEError(err):
				// Respond as per RFC 7230 Section 3.3.1 which says,
				//      A server that receives a request message with a
				//      transfer coding it does not understand SHOULD
				//      respond with 501 (Unimplemented).
				code := StatusNotImplemented

				// We purposefully aren't echoing back the transfer-encoding's value,
				// so as to mitigate the risk of cross side scripting by an attacker.
				fmt.Fprintf(c.rwc, "HTTP/1.1 %d %s%sUnsupported transfer encoding", code, StatusText(code), errorHeaders)
				return

			case isCommonNetReadError(err):
				return // don't reply

			default:
				if v, ok := err.(statusError); ok {
					fmt.Fprintf(c.rwc, "HTTP/1.1 %d %s: %s%s%d %s: %s", v.code, StatusText(v.code), v.text, errorHeaders, v.code, StatusText(v.code), v.text)
					return
				}
				publicErr := "400 Bad Request"
				fmt.Fprintf(c.rwc, "HTTP/1.1 "+publicErr+errorHeaders+publicErr)
				return
			}
		}

		// Expect 100 Continue support
		req := w.req
		if req.expectsContinue() {
			if req.ProtoAtLeast(1, 1) && req.ContentLength != 0 {
				// Wrap the Body reader with one that replies on the connection
				req.Body = &expectContinueReader{readCloser: req.Body, resp: w}
				w.canWriteContinue.setTrue()
			}
		} else if req.Header.get("Expect") != "" {
			w.sendExpectationFailed()
			return
		}

		c.curReq.Store(w)

		if requestBodyRemains(req.Body) {
			registerOnHitEOF(req.Body, w.conn.r.startBackgroundRead)
		} else {
			w.conn.r.startBackgroundRead()
		}

		// HTTP cannot have multiple simultaneous active requests.[*]
		// Until the server replies to this request, it can't read another,
		// so we might as well run the handler in this goroutine.
		// [*] Not strictly true: HTTP pipelining. We could let them all process
		// in parallel even if their responses need to be serialized.
		// But we're not going to implement HTTP pipelining because it
		// was never deployed in the wild and the answer is HTTP/2.
		serverHandler{c.server}.ServeHTTP(w, w.req)
		w.cancelCtx()
		if c.hijacked() {
			return
		}
		w.finishRequest()
		if !w.shouldReuseConnection() {
			if w.requestBodyLimitHit || w.closedRequestBodyEarly() {
				c.closeWriteAndWait()
			}
			return
		}
		c.setState(c.rwc, StateIdle, runHooks)
		c.curReq.Store((*response)(nil))

		if !w.conn.server.doKeepAlives() {
			// We're in shutdown mode. We might've replied
			// to the user without "Connection: close" and
			// they might think they can send another
			// request, but such is life with HTTP/1.1.
			return
		}

		if d := c.server.idleTimeout(); d != 0 {
			c.rwc.SetReadDeadline(time.Now().Add(d))
			if _, err := c.bufr.Peek(4); err != nil {
				return
			}
		}
		c.rwc.SetReadDeadline(time.Time{})
	}
}
```


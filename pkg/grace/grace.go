// Package grace
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:59
package grace

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	// PreSignal 在信号之前添加过滤器
	PreSignal = iota
	// PostSignal 在信号之后添加过滤器
	PostSignal
	// StateInit 表示应用程序正在初始化
	StateInit
	// StateRunning 表示应用程序正在运行
	StateRunning
	// StateShuttingDown 表示应用程序正在关闭
	StateShuttingDown
	// StateTerminate 表示应用程序已经关闭
	StateTerminate
)

var (
	regLock              *sync.Mutex //
	runningServers       map[string]*Server
	runningServersOrder  []string
	socketPtrOffsetMap   map[string]uint
	runningServersForked bool

	// DefaultReadTimeOut 读取超时时间
	DefaultReadTimeOut time.Duration
	// DefaultWriteTimeOut 写入超时时间
	DefaultWriteTimeOut time.Duration
	// DefaultMaxHeaderBytes http标头的最大大小，默认是0，没有限制
	DefaultMaxHeaderBytes int
	// DefaultTimeout 关闭服务器的超时时间，默认是60s
	DefaultTimeout = 60 * time.Second

	isChild     bool
	socketOrder string

	// hookableSignals 拦截的信号
	hookableSignals []os.Signal
)

func init() {
	//flag.BoolVar(&isChild, "graceful", false, "listen on open fd (after forking)")
	//flag.StringVar(&socketOrder, "socketorder", "", "previous initialization order - used when more than one listener was started")
	socketOrder = os.Getenv("ENDLESS_SOCKET_ORDER")
	isChild = os.Getenv("ENDLESS_CONTINUE") != ""

	regLock = &sync.Mutex{}
	runningServers = make(map[string]*Server)
	runningServersOrder = []string{}
	socketPtrOffsetMap = make(map[string]uint)

	hookableSignals = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	}
}

// ServerOption 连接配置
type ServerOption func(*Server)

// WithShutdownCallback 设置关闭回调
func WithShutdownCallback(shutdownCallback func()) ServerOption {
	return func(srv *Server) {
		srv.shutdownCallbacks = append(srv.shutdownCallbacks, shutdownCallback)
	}
}

// NewServer 启动一个服务
func NewServer(addr string, handler http.Handler, opts ...ServerOption) (srv *Server) {
	regLock.Lock()
	defer regLock.Unlock()

	//if !flag.Parsed() {
	//	flag.Parse()
	//}
	if len(socketOrder) > 0 {
		for i, addr := range strings.Split(socketOrder, ",") {
			socketPtrOffsetMap[addr] = uint(i)
		}
	} else {
		socketPtrOffsetMap[addr] = uint(len(runningServersOrder))
	}

	srv = &Server{
		sigChan: make(chan os.Signal),
		isChild: isChild,
		SignalHooks: map[int]map[os.Signal][]func(){
			PreSignal: {
				syscall.SIGHUP:  {},
				syscall.SIGINT:  {},
				syscall.SIGTERM: {},
			},
			PostSignal: {
				syscall.SIGHUP:  {},
				syscall.SIGINT:  {},
				syscall.SIGTERM: {},
			},
		},
		state:        StateInit,
		Network:      "tcp",
		terminalChan: make(chan error), // no cache channel
	}
	srv.Server = &http.Server{
		Addr:           addr,
		ReadTimeout:    DefaultReadTimeOut,
		WriteTimeout:   DefaultWriteTimeOut,
		MaxHeaderBytes: DefaultMaxHeaderBytes,
		Handler:        handler,
	}

	for _, opt := range opts {
		opt(srv)
	}

	runningServersOrder = append(runningServersOrder, addr)
	runningServers[addr] = srv
	return srv
}

// ListenAndServe 参考 http.ListenAndServe
func ListenAndServe(addr string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServe()
}

// ListenAndServeTLS 参考 http.ListenAndServeTLS
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServeTLS(certFile, keyFile)
}

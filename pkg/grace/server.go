// Package grace
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-09-11 17:22
package grace

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/cloudflare/tableflip"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Server 嵌入 http.Server
type Server struct {
	*http.Server
	ln                net.Listener
	SignalHooks       map[int]map[os.Signal][]func()
	sigChan           chan os.Signal
	isChild           bool
	state             uint8
	Network           string
	terminalChan      chan error
	shutdownCallbacks []func()
	upg               *tableflip.Upgrader
}

//func init() {
//	upg, err := tableflip.New(tableflip.Options{
//		PIDFile: ".pid",
//	})
//	if err != nil {
//		panic(err)
//	}
//	//defer func() {
//	//	fmt.Println("upg.Stop()")
//	//	upg.Stop()
//	//}()
//}

// Serve accepts incoming connections on the Listener l
// and creates a new service goroutine for each.
// The service goroutines read requests and then call srv.Handler to reply to them.
func (srv *Server) Serve() (err error) {
	return srv.internalServe(srv.ln)
}

func (srv *Server) ServeWithListener(ln net.Listener) (err error) {
	srv.ln = ln
	go srv.handleSignals()
	return srv.internalServe(ln)
}

func (srv *Server) internalServe(ln net.Listener) (err error) {
	srv.state = StateRunning
	defer func() { srv.state = StateTerminate }()
	defer func() {
		srv.upg.Stop()
		fmt.Println("upg.Stop()")
	}()

	go func() {
		// When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS
		// immediately return ErrServerClosed. Make sure the program doesn't exit
		// and waits instead for Shutdown to return.
		if err = srv.Server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(syscall.Getpid(), "Server.Serve() error:", err)
		}
	}()

	if err = srv.upg.Ready(); err != nil {
		panic(err)
	}

	// wait for Shutdown to return
	if shutdownErr := <-srv.terminalChan; shutdownErr != nil {
		return shutdownErr
	}

	log.Println(syscall.Getpid(), ln.Addr(), "Listener closed.")

	time.AfterFunc(30*time.Second, func() {
		os.Exit(1)
	})

	//srv.upg.Shutdown(context.Background())
	return
}

// ListenAndServe listens on the TCP network address srv.Addr and then calls Serve
// to handle requests on incoming connections. If srv.Addr is blank, ":http" is
// used.
func (srv *Server) ListenAndServe() (err error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}

	go srv.handleSignals()

	srv.ln, err = srv.getListener(addr)
	if err != nil {
		log.Println(err)
		return err
	}

	if srv.isChild {
		process, err := os.FindProcess(os.Getppid())
		if err != nil {
			log.Println(err)
			return err
		}
		err = process.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}
	}

	log.Println(os.Getpid(), srv.Addr)
	return srv.Serve()
}

// ListenAndServeTLS listens on the TCP network address srv.Addr and then calls
// Serve to handle requests on incoming TLS connections.
//
// Filenames containing a certificate and matching private key for the server must
// be provided. If the certificate is signed by a certificate authority, the
// certFile should be the concatenation of the server's certificate followed by the
// CA's certificate.
//
// If srv.Addr is blank, ":https" is used.
func (srv *Server) ListenAndServeTLS(certFile, keyFile string) (err error) {
	ln, err := srv.ListenTLS(certFile, keyFile)
	if err != nil {
		return err
	}

	return srv.ServeTLS(ln)
}

func (srv *Server) ListenTLS(certFile string, keyFile string) (net.Listener, error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	if srv.TLSConfig == nil {
		srv.TLSConfig = &tls.Config{}
	}
	if srv.TLSConfig.NextProtos == nil {
		srv.TLSConfig.NextProtos = []string{"http/1.1"}
	}

	srv.TLSConfig.Certificates = make([]tls.Certificate, 1)
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	srv.TLSConfig.Certificates[0] = cert

	go srv.handleSignals()

	ln, err := srv.getListener(addr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, srv.TLSConfig)
	return tlsListener, nil
}

// ListenAndServeMutualTLS listens on the TCP network address srv.Addr and then calls
// Serve to handle requests on incoming mutual TLS connections.
func (srv *Server) ListenAndServeMutualTLS(certFile, keyFile, trustFile string) (err error) {
	ln, err := srv.ListenMutualTLS(certFile, keyFile, trustFile)
	if err != nil {
		return err
	}

	return srv.ServeTLS(ln)
}

func (srv *Server) ServeTLS(ln net.Listener) error {
	if srv.isChild {
		process, err := os.FindProcess(os.Getppid())
		if err != nil {
			log.Println(err)
			return err
		}
		err = process.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}
	}

	go srv.handleSignals()
	return srv.internalServe(ln)
}

func (srv *Server) ListenMutualTLS(certFile string, keyFile string, trustFile string) (net.Listener, error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	if srv.TLSConfig == nil {
		srv.TLSConfig = &tls.Config{}
	}
	if srv.TLSConfig.NextProtos == nil {
		srv.TLSConfig.NextProtos = []string{"http/1.1"}
	}

	srv.TLSConfig.Certificates = make([]tls.Certificate, 1)
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	srv.TLSConfig.Certificates[0] = cert
	srv.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	pool := x509.NewCertPool()
	data, err := os.ReadFile(trustFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	pool.AppendCertsFromPEM(data)
	srv.TLSConfig.ClientCAs = pool
	log.Println("Mutual HTTPS")
	go srv.handleSignals()

	ln, err := srv.getListener(addr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, srv.TLSConfig)
	return tlsListener, nil
}

// getListener either opens a new socket to listen on, or takes the acceptor socket
// it got passed when restarted.
func (srv *Server) getListener(laddr string) (l net.Listener, err error) {
	if srv.isChild {
		var ptrOffset uint
		if len(socketPtrOffsetMap) > 0 {
			ptrOffset = socketPtrOffsetMap[laddr]
			log.Println("laddr", laddr, "ptr offset", socketPtrOffsetMap[laddr])
		}

		f := os.NewFile(uintptr(3+ptrOffset), "")
		l, err = net.FileListener(f)
		if err != nil {
			err = fmt.Errorf("net.FileListener error: %v", err)
			return
		}
	} else {
		//l, err = net.Listen(srv.Network, laddr)
		l, err = srv.upg.Fds.Listen(srv.Network, laddr)
		if err != nil {
			err = fmt.Errorf("net.Listen error: %v", err)
			return
		}
	}
	return
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// handleSignals listens for os Signals and calls any hooked in function that the
// user had registered with the signal.
func (srv *Server) handleSignals() {
	var sig os.Signal

	signal.Notify(
		srv.sigChan,
		hookableSignals...,
	)

	pid := syscall.Getpid()
	for {
		sig = <-srv.sigChan
		srv.signalHooks(PreSignal, sig)
		switch sig {
		case syscall.SIGHUP:
			log.Println(pid, "Received SIGHUP. forking.")
			err := srv.upg.Upgrade()
			if err != nil {
				log.Println("Upgrade failed:", err)
				continue
			}
		case syscall.SIGINT:
			log.Println(pid, "Received SIGINT.")
			srv.upg.Stop()
			srv.shutdown()
		case syscall.SIGTERM:
			log.Println(pid, "Received SIGTERM.")
			srv.upg.Stop()
			srv.shutdown()
		default:
			log.Printf("Received %v: nothing i care about...\n", sig)
		}
		srv.signalHooks(PostSignal, sig)
	}
}

func (srv *Server) signalHooks(ppFlag int, sig os.Signal) {
	if _, notSet := srv.SignalHooks[ppFlag][sig]; !notSet {
		return
	}
	for _, f := range srv.SignalHooks[ppFlag][sig] {
		f()
	}
}

// shutdown closes the listener so that no new connections are accepted. it also
// starts a goroutine that will serverTimeout (stop all running requests) the server
// after DefaultTimeout.
func (srv *Server) shutdown() {
	if srv.state != StateRunning {
		return
	}

	srv.state = StateShuttingDown
	log.Println(syscall.Getpid(), "Waiting for connections to finish...")
	ctx := context.Background()
	if DefaultTimeout >= 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()
	}
	for _, shutdownCallback := range srv.shutdownCallbacks {
		shutdownCallback()
	}
	srv.terminalChan <- srv.Server.Shutdown(ctx)
}

func (srv *Server) fork() (err error) {
	regLock.Lock()
	defer regLock.Unlock()
	if runningServersForked {
		return
	}
	runningServersForked = true

	files := make([]*os.File, len(runningServers))
	orderArgs := make([]string, len(runningServers))
	for _, srvPtr := range runningServers {
		f, _ := srvPtr.ln.(*net.TCPListener).File()
		files[socketPtrOffsetMap[srvPtr.Server.Addr]] = f
		orderArgs[socketPtrOffsetMap[srvPtr.Server.Addr]] = srvPtr.Server.Addr
	}

	log.Println(files)
	path := os.Args[0]
	var args []string
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if arg == "-graceful" {
				break
			}
			args = append(args, arg)
		}
	}
	args = append(args, "-graceful")
	if len(runningServers) > 1 {
		args = append(args, fmt.Sprintf(`-socketorder=%s`, strings.Join(orderArgs, ",")))
		log.Println(args)
	}
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = files
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Restart: Failed to launch, error: %v", err)
	}
	log.Println(syscall.Getpid(), "Received SIGHUP. forking1111.")

	return
}

// RegisterSignalHook registers a function to be run PreSignal or PostSignal for a given signal.
func (srv *Server) RegisterSignalHook(ppFlag int, sig os.Signal, f func()) (err error) {
	if ppFlag != PreSignal && ppFlag != PostSignal {
		err = fmt.Errorf("invalid ppFlag argument. Must be either grace.PreSignal or grace.PostSignal")
		return
	}
	for _, s := range hookableSignals {
		if s == sig {
			srv.SignalHooks[ppFlag][sig] = append(srv.SignalHooks[ppFlag][sig], f)
			return
		}
	}
	err = fmt.Errorf("signal '%v' is not supported", sig)
	return
}

package main

// This is a launcher for a generic zcl language server, which can work with
// basic zcl syntax but doesn't have any knowledge of which blocks and
// attributes are valid for a particular application.
//
// This can be used directly if desired, particularly for simple applications,
// but non-trivial applications ought to provide their own version of this
// which configures the underlying library with application-specific
// information.
//
// This is based on sourcegraph's language server for Go.

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/sourcegraph/jsonrpc2"
	"github.com/zclconf/zcl-language-server/zclsrv"

	_ "net/http/pprof"
)

var (
	mode         = flag.String("mode", "stdio", "communication mode (stdio|tcp)")
	addr         = flag.String("addr", ":4389", "server listen address (tcp)")
	trace        = flag.Bool("trace", false, "print all requests and responses")
	logfile      = flag.String("logfile", "", "also log to this file (in addition to stderr)")
	printVersion = flag.Bool("version", false, "print version and exit")
	pprof        = flag.String("pprof", ":6060", "start a pprof http server (https://golang.org/pkg/net/http/pprof/)")
	freeosmemory = flag.Bool("freeosmemory", true, "aggressively free memory back to the OS")
)

const version = "v0-dev"

func main() {
	flag.Parse()
	log.SetFlags(0)

	// Start pprof server, if desired.
	if *pprof != "" {
		go func() {
			log.Println(http.ListenAndServe(*pprof, nil))
		}()
	}

	if *freeosmemory {
		go freeOSMemory()
	}

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if *printVersion {
		fmt.Println(version)
		return nil
	}

	var logW io.Writer
	if *logfile == "" {
		logW = os.Stderr
	} else {
		f, err := os.Create(*logfile)
		if err != nil {
			return err
		}
		defer f.Close()
		logW = io.MultiWriter(os.Stderr, f)
	}
	log.SetOutput(logW)

	var connOpt []jsonrpc2.ConnOpt
	if *trace {
		connOpt = append(connOpt, jsonrpc2.LogMessages(log.New(logW, "", 0)))
	}

	switch *mode {
	case "tcp":
		lis, err := net.Listen("tcp", *addr)
		if err != nil {
			return err
		}
		defer lis.Close()

		log.Println("zcl-language-server: listening on", *addr)
		for {
			conn, err := lis.Accept()
			if err != nil {
				return err
			}
			jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}), zclsrv.NewHandler(), connOpt...)
		}

	case "stdio":
		log.Println("zcl-language-server: reading on stdin, writing on stdout")
		<-jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{}), zclsrv.NewHandler(), connOpt...).DisconnectNotify()
		log.Println("connection closed")
		return nil

	default:
		return fmt.Errorf("invalid mode %q", *mode)
	}
}

type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}

func freeOSMemory() {
	for {
		time.Sleep(1 * time.Second)
		debug.FreeOSMemory()
	}
}

package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/cbandy/pgtwixt"
	"github.com/go-kit/kit/log"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	listen, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}

	go func() {
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt)
		for s := range signals {
			fmt.Printf("Got signal %v\n", s)
			os.Exit(1)
		}
	}()

	var connector pgtwixt.Connector

	if strings.HasPrefix(os.Args[2], "/") {
		dialer := pgtwixt.UnixDialer{
			Addr:  os.Args[2],
			Debug: logger.Log,
		}
		connector.Dial = dialer.Dial
	} else {
		dialer := pgtwixt.TCPDialer{
			Addr:    os.Args[2],
			Debug:   logger.Log,
			SSLMode: "prefer",
			SSLConfig: tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS12,
				Renegotiation:      tls.RenegotiateFreelyAsClient,
			},
		}
		connector.Dial = dialer.Dial
	}

	proxy := pgtwixt.Proxy{Info: logger.Log, Startup: connector.Startup}

	srv := pgtwixt.Server{
		Debug: logger.Log,
		Info:  logger.Log,

		Cancel: func(c pgtwixt.CancellationKey) {
			if err := connector.Cancel(c); err != nil {
				logger.Log("msg", "Error during cancel", "error", err)
			}
		},
		Session: func(fe pgtwixt.FrontendStream, startup map[string]string) {
			fmt.Printf("%#v\n", startup)
			proxy.Run(fe, startup)
		},
	}

	err = srv.Serve(listen)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/cbandy/pgtwixt"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	listen, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}

	go func() {
		metrics := &http.Server{
			Addr:         os.Args[2],
			ReadTimeout:  4 * time.Second,
			WriteTimeout: 4 * time.Second,
			Handler: promhttp.InstrumentMetricHandler(
				metricRegistry, promhttp.HandlerFor(
					metricGatherer, promhttp.HandlerOpts{},
				),
			),
		}

		err := metrics.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt)
		for s := range signals {
			fmt.Printf("Got signal %v\n", s)
			os.Exit(1)
		}
	}()

	var connector pgtwixt.Connector

	if strings.HasPrefix(os.Args[3], "/") {
		connector.Dialer = pgtwixt.UnixDialer{
			Addr:  os.Args[3],
			Debug: logger.Log,
		}
	} else {
		connector.Dialer = pgtwixt.TCPDialer{
			Addr:    os.Args[3],
			Debug:   logger.Log,
			SSLMode: "prefer",
			SSLConfig: tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS12,
				Renegotiation:      tls.RenegotiateFreelyAsClient,
			},
		}
	}

	proxy := pgtwixt.Proxy{
		Info: logger.Log,

		Startup: connector.Startup,

		CountConnect:    metrics.backend.connects.With(prometheus.Labels{"backend": "yes", "host": os.Args[3]}).Inc,
		CountDisconnect: metrics.backend.disconnects.With(prometheus.Labels{"backend": "yes", "host": os.Args[3]}).Inc,
	}

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

		CountConnect:    metrics.frontend.connects.With(prometheus.Labels{"frontend": "yes", "bind": listen.Addr().String()}).Inc,
		CountDisconnect: metrics.frontend.disconnects.With(prometheus.Labels{"frontend": "yes", "bind": listen.Addr().String()}).Inc,
	}

	err = srv.Serve(listen)
	if err != nil {
		panic(err)
	}
}

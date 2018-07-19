package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
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

	var connstr pgtwixt.ConnectionString
	err = connstr.Parse(os.Args[3])
	if err != nil {
		panic(err)
	}

	ds, err := Connector{Debug: logger.Log}.Dialers(connstr)
	if err != nil {
		panic(err)
	}

	connector := pgtwixt.Connector{Dialer: ds[0]}
	proxy := pgtwixt.Proxy{
		Info: logger.Log,

		Startup: connector.Startup,

		CountConnect: func() func() {
			var (
				connections = metrics.backend.connections.With(prometheus.Labels{"backend": "yes", "host": os.Args[3]})
				connects    = metrics.backend.connects.With(prometheus.Labels{"backend": "yes", "host": os.Args[3]})
			)
			return func() { connections.Inc(); connects.Inc() }
		}(),
		CountDisconnect: func() func() {
			var (
				connections = metrics.backend.connections.With(prometheus.Labels{"backend": "yes", "host": os.Args[3]})
				disconnects = metrics.backend.disconnects.With(prometheus.Labels{"backend": "yes", "host": os.Args[3]})
			)
			return func() { connections.Dec(); disconnects.Inc() }
		}(),
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

		CountConnect: func() func() {
			var (
				connections = metrics.frontend.connections.With(prometheus.Labels{"frontend": "yes", "bind": listen.Addr().String()})
				connects    = metrics.frontend.connects.With(prometheus.Labels{"frontend": "yes", "bind": listen.Addr().String()})
			)
			return func() { connections.Inc(); connects.Inc() }
		}(),
		CountDisconnect: func() func() {
			var (
				connections = metrics.frontend.connections.With(prometheus.Labels{"frontend": "yes", "bind": listen.Addr().String()})
				disconnects = metrics.frontend.disconnects.With(prometheus.Labels{"frontend": "yes", "bind": listen.Addr().String()})
			)
			return func() { connections.Dec(); disconnects.Inc() }
		}(),
	}

	err = srv.Serve(listen)
	if err != nil {
		panic(err)
	}
}

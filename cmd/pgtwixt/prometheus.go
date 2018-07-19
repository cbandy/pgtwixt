package main

import "github.com/prometheus/client_golang/prometheus"

var metricGatherer prometheus.Gatherers
var metricRegistry *prometheus.Registry
var metrics struct {
	backend struct {
		connections *prometheus.GaugeVec
		connects    *prometheus.CounterVec
		disconnects *prometheus.CounterVec
	}
	frontend struct {
		connections *prometheus.GaugeVec
		connects    *prometheus.CounterVec
		disconnects *prometheus.CounterVec
	}
}

func init() {
	metricRegistry = prometheus.NewRegistry()
	metricRegistry.MustRegister(prometheus.NewGoCollector())
	metricGatherer = append(metricGatherer, metricRegistry)

	metricGatherer = append(metricGatherer, func() (backend *prometheus.Registry) {
		backend = prometheus.NewPedanticRegistry()

		metrics.backend.connections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "pgtwixt_connections",
			Help: "Current number of connections to both frontends and backends.",
		}, []string{"host"})
		backend.MustRegister(metrics.backend.connections)

		metrics.backend.connects = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "pgtwixt_connects_total",
			Help: "Total number of connects to both frontends and backends.",
		}, []string{"host"})
		backend.MustRegister(metrics.backend.connects)

		metrics.backend.disconnects = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "pgtwixt_disconnects_total",
			Help: "Total number of disconnects from both frontends and backends.",
		}, []string{"host"})
		backend.MustRegister(metrics.backend.disconnects)

		return
	}())

	metricGatherer = append(metricGatherer, func() (frontend *prometheus.Registry) {
		frontend = prometheus.NewPedanticRegistry()

		metrics.frontend.connections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "pgtwixt_connections",
			Help: "Current number of connections to both frontends and backends.",
		}, []string{"bind"})
		frontend.MustRegister(metrics.frontend.connections)

		metrics.frontend.connects = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "pgtwixt_connects_total",
			Help: "Total number of connects to both frontends and backends.",
		}, []string{"bind"})
		frontend.MustRegister(metrics.frontend.connects)

		metrics.frontend.disconnects = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "pgtwixt_disconnects_total",
			Help: "Total number of disconnects from both frontends and backends.",
		}, []string{"bind"})
		frontend.MustRegister(metrics.frontend.disconnects)

		return
	}())
}

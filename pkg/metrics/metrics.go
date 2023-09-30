// Package metrics provides instrumentation for monitoring and alerting
// using Prometheus metrics in the hpxd application.
//
// Author: zakaria.elbouwab
package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// GitPullCounter tracks the number of times the config is pulled from Git.
	//
	// This metric is a counter that can be labeled with 'status' which can be either
	// 'success' or 'failure' depending on the outcome of the pull operation.
	GitPullCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hpxd_git_pulls_total",
			Help: "Total number of times the config is pulled from Git",
		},
		[]string{"status"}, // success or failure
	)

	// HaproxyReloadCounter tracks the number of times HAProxy is reloaded.
	//
	// This is a simple counter metric without labels. It increments every time HAProxy
	// is reloaded by the hpxd application.
	HaproxyReloadCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hpxd_haproxy_reloads_total",
			Help: "Total number of times HAProxy is reloaded",
		},
	)

	// InvalidConfigCounter tracks the number of times an invalid config is detected.
	//
	// This counter metric increments each time the hpxd application detects
	// an invalid configuration for HAProxy.
	InvalidConfigCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hpxd_invalid_configs_total",
			Help: "Total number of times an invalid config is detected",
		},
	)

	// ApplicationInfo provides details about the running application.
	//
	// This gauge metric is labeled with 'version', 'commit', and 'buildDate' to
	// give insights into which version of the application is currently running
	// and its associated metadata.
	ApplicationInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "application_info",
			Help: "Application info (version, commit, buildDate)",
		},
		[]string{"version", "commit", "buildDate"},
	)
)

func init() {
	// Registering the metrics with Prometheus's default registry ensures they are
	// exposed for scraping by a Prometheus server.
	prometheus.MustRegister(GitPullCounter, HaproxyReloadCounter, InvalidConfigCounter, ApplicationInfo)
}

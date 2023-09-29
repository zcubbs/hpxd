package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	GitPullCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hpxd_git_pulls_total",
			Help: "Total number of times the config is pulled from Git",
		},
		[]string{"status"}, // success or failure
	)

	HaproxyReloadCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hpxd_haproxy_reloads_total",
			Help: "Total number of times HAProxy is reloaded",
		},
	)

	InvalidConfigCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hpxd_invalid_configs_total",
			Help: "Total number of times an invalid config is detected",
		},
	)
)

func init() {
	// Register the metrics with Prometheus's default registry
	prometheus.MustRegister(GitPullCounter, HaproxyReloadCounter, InvalidConfigCounter)
}

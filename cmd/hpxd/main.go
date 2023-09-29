package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zcubbs/hpxd/pkg/git"
	"github.com/zcubbs/hpxd/pkg/haproxy"
	"github.com/zcubbs/hpxd/pkg/metrics"
)

const (
	defaultPollingInterval = 5 * time.Second
	prometheusDefaultPort  = 9100
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

type Configuration struct {
	RepoURL           string        `mapstructure:"repoURL"`
	Branch            string        `mapstructure:"branch"`
	HaproxyConfigPath string        `mapstructure:"haproxyConfigPath"`
	PollingInterval   time.Duration `mapstructure:"pollingInterval"`
	EnablePrometheus  bool          `mapstructure:"enablePrometheus"`
	PrometheusPort    int           `mapstructure:"prometheusPort"`

	Version string
	Commit  string
	Date    string
}

// setupConfig reads in config file.
// example: `./hpxd -config=/path/to/configs`
func setupConfig() *Configuration {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs", "Path to configuration directory")
	flag.Parse()

	viper.SetConfigName("hpxd")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	viper.SetDefault("enablePrometheus", false)
	viper.SetDefault("prometheusPort", prometheusDefaultPort)
	viper.SetDefault("pollingInterval", defaultPollingInterval)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file, %s", err)
	}

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		logrus.Fatalf("Unable to unmarshal into struct, %v", err)
	}

	config.Version = Version
	config.Commit = Commit
	config.Date = Date

	return &config
}

func validateConfig(config *Configuration) error {
	if config.RepoURL == "" {
		return errors.New("missing required config: repoURL")
	}

	if config.Branch == "" {
		return errors.New("missing required config: branch")
	}

	if config.HaproxyConfigPath == "" {
		return errors.New("missing required config: haproxyConfigPath")
	}
	return nil
}

func startMetricsEndpoint(port int) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			logrus.Fatalf("Error starting metrics endpoint: %v", err)
		}
	}()
}

func main() {
	config := setupConfig()
	if err := validateConfig(config); err != nil {
		logrus.Fatal(err)
	}

	gitHandler := git.NewHandler(config.RepoURL, config.Branch, config.HaproxyConfigPath)
	haproxyHandler := haproxy.NewHandler(config.HaproxyConfigPath)

	if config.EnablePrometheus {
		startMetricsEndpoint(config.PrometheusPort)
	}

	for {
		updated, err := gitHandler.PullAndUpdate()
		if err != nil {
			logrus.Errorf("Error while pulling updates: %v", err)
			// Update Prometheus metric for failed Git pull
			metrics.GitPullCounter.WithLabelValues("failure").Inc()
			time.Sleep(config.PollingInterval)
			continue
		}

		if updated {
			// Update Prometheus metric for successful Git pull
			metrics.GitPullCounter.WithLabelValues("success").Inc()

			// Check if new configuration is valid
			if err := haproxyHandler.ValidateConfig(); err != nil {
				logrus.Errorf("Pulled HAProxy configuration is invalid: %v", err)
				// Update Prometheus metric for invalid config
				metrics.InvalidConfigCounter.Inc()
			} else {
				// Handle HAProxy reload logic
				if err := haproxyHandler.Reload(); err != nil {
					logrus.Errorf("Failed to reload HAProxy: %v", err)
				} else {
					// Update Prometheus metric for successful HAProxy reload
					metrics.HaproxyReloadCounter.Inc()
					logrus.Info("Configuration updated and HAProxy reloaded successfully!")
				}
			}
		}
		time.Sleep(config.PollingInterval)
	}
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	RepoURL string `mapstructure:"repoURL"`
	Branch  string `mapstructure:"branch"`
	Path    string `mapstructure:"path"`

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

	if config.Path == "" {
		return errors.New("missing required config: path")
	}

	if config.HaproxyConfigPath == "" {
		return errors.New("missing required config: haproxyConfigPath")
	}
	return nil
}

func startMetricsEndpoint(port int) {
	// register app version info
	metrics.ApplicationInfo.WithLabelValues(Version, Commit, Date).Set(1)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		server := &http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			ReadHeaderTimeout: 3 * time.Second,
		}
		err := server.ListenAndServe()
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

	logrus.Infof("Starting hpxd version %s (%s) built on %s",
		config.Version,
		config.Commit,
		config.Date,
	)

	gitHandler := git.NewHandler(config.RepoURL, config.Branch, config.Path, config.HaproxyConfigPath)
	haproxyHandler := haproxy.NewHandler(config.HaproxyConfigPath)

	if config.EnablePrometheus {
		startMetricsEndpoint(config.PrometheusPort)
	}

	update(gitHandler, haproxyHandler, config)
}

// update is the main loop of hpxd. This is what happens in the loop:
//
// 1. HAProxy's configuration is fetched from git.
//
// 2. The fetched configuration is validated. If it's invalid, the loop continues.
//
// 3. If the configuration is valid, it's applied and HAProxy is reloaded.
func update(gitHandler *git.Handler, haproxyHandler *haproxy.Handler, config *Configuration) {
	for {
		configPath, updated, err := gitHandler.PullAndUpdate()
		if err != nil {
			logrus.Errorf("Error while pulling updates: %v", err)
			// Update Prometheus metric for failed Git pull
			metrics.GitPullCounter.WithLabelValues("failure").Inc()
			time.Sleep(config.PollingInterval)
			continue
		}

		if updated {
			// Temporarily create a handler for validation
			tempHandler := haproxy.NewHandler(configPath)

			// Check if new configuration is valid
			if err := tempHandler.ValidateConfig(); err != nil {
				logrus.Errorf("Pulled HAProxy configuration is invalid: %v", err)
				// Update Prometheus metric for invalid config
				metrics.InvalidConfigCounter.Inc()
			} else {
				// If valid, update the actual config and reload HAProxy
				copyConfig(configPath, config.HaproxyConfigPath)

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

func copyConfig(src, dest string) {
	input, err := os.ReadFile(filepath.Clean(src))
	if err != nil {
		logrus.Errorf("Failed to read config from source: %v", err)
		return
	}

	err = os.WriteFile(dest, input, 0600)
	if err != nil {
		logrus.Errorf("Error writing config to destination: %v", err)
	}
}

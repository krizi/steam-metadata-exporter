package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"steam-exporter/config"
	"steam-exporter/feature"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
)

func main() {
	logger := logrus.New()
	cfg := config.LoadConfig(logger.WithField("lib", "sirupsen/logrus"))

	// Initialize logger

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.Fatalf("Invalid log level: %v", err)
	}
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})
	appLogger := logger.WithField("app", "steam-exporter")

	scheduler := gocron.NewScheduler(time.UTC)

	for featureName, featureConfig := range cfg.Features {
		if featureConfig.Enabled {
			var featureLogger = appLogger.WithField("task", featureName)
			instance, _ := feature.CreateInstance(featureName, featureLogger, cfg, featureConfig)
			initError := instance.InitializeFeature()
			if initError != nil {
				featureLogger.Errorf("Initialization failed: %v", initError)
			}

			_, _ = scheduler.Every(featureConfig.Schedule).Do(func(feature feature.Feature) {
				if feature.Execute() != nil {
					featureLogger.Errorf("Execution failed: %v", err)
				}
			}, instance)
		}
	}

	appLogger.Info("Starting Steam Exporter...")
	scheduler.StartAsync()

	r := gin.New()
	r.Use(ginlogrus.Logger(appLogger.WithField("lib", "gin-gonic/gin")), gin.Recovery())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	appLogger.Infof("Starting server on %s", cfg.MetricsPort)
	if r.Run(cfg.MetricsPort) != nil {
		appLogger.Fatalf("Failed to start server: %v", err)
	}
}

package feature

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"steam-exporter/config"
)

type Feature interface {
	InitializeFeature() error
	Execute() error
}

func CreateInstance(name string,
	logger *logrus.Entry,
	cfg *config.SteamConfig,
	featureConfig config.FeatureConfig) (Feature, error) {
	switch name {
	case "owned_games":
		return &GamesTotalFeature{
			logger,
			cfg,
			featureConfig,
		}, nil
	case "playtime":
		return &PlaytimeFeature{
			logger,
			cfg,
			featureConfig,
		}, nil
	}

	return nil, fmt.Errorf("unknown feature: %s", name)
}

package feature

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"steam-exporter/config"
)

type GamesTotalFeature struct {
	logger        *logrus.Entry
	cfg           *config.SteamConfig
	featureConfig config.FeatureConfig
}

func (o *GamesTotalFeature) InitializeFeature() error {
	o.logger.Info("Initializing GamesTotalFeature...")
	prometheus.MustRegister(gamesTotalMetric)
	return nil
}

var gamesTotalMetric = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "steam_games_total",
		Help: "Total number of owned games",
	},
)

func (o *GamesTotalFeature) Execute() error {
	o.logger.Infof("Executing GamesTotalFeature for user %s", o.cfg.UserID)
	client := resty.New()
	response, err := client.R().
		SetQueryParams(map[string]string{
			"key":                       o.cfg.APIKey,
			"steamid":                   o.cfg.UserID,
			"include_played_free_games": "true",
		}).
		SetHeader("Accept", "application/json").
		Get(o.featureConfig.APIURL)
	if err != nil {
		return err
	}

	var result struct {
		Response struct {
			GameCount float64 `json:"game_count"`
		} `json:"response"`
	}

	if json.Unmarshal(response.Body(), &result) != nil {
		return err
	}

	gamesTotalMetric.Set(result.Response.GameCount)
	return nil
}

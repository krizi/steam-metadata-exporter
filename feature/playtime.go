package feature

import (
	"fmt"

	"steam-exporter/config"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type PlaytimeFeature struct {
	logger        *logrus.Entry
	cfg           *config.SteamConfig
	featureConfig config.FeatureConfig
}

type APIResult struct {
	Response struct {
		Games []struct {
			AppID      int     `json:"appid"`
			Name       string  `json:"name"`
			Playtime   float64 `json:"playtime_forever"`
			ImgIconURL string  `json:"img_icon_url"`
			ImgLogoURL string  `json:"img_logo_url"`
		} `json:"games"`
	} `json:"response"`
}

var steamPlaytimeMinutesTotal = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "steam_playtime_minutes_total",
		Help: "Total playtime in minutes for each game on a specific platform",
	},
	[]string{"appid", "name", "img_icon_url", "http_img_icon_url", "img_logo_url"},
)

func (playtime *PlaytimeFeature) InitializeFeature() error {
	playtime.logger.Info("Initializing GamesTotalFeature...")
	prometheus.MustRegister(steamPlaytimeMinutesTotal)
	return nil
}

func (playtime *PlaytimeFeature) Execute() error {
	playtime.logger.Infof("Executing PlaytimeFeature for user %s", playtime.cfg.UserID)
	client := resty.New()

	var result APIResult

	_, err := client.R().
		SetQueryParams(map[string]string{
			"key":                       playtime.cfg.APIKey,
			"steamid":                   playtime.cfg.UserID,
			"include_played_free_games": "true",
			"include_appinfo":           "true",
			"include_extended_appinfo":  "true",
		}).
		SetHeader("Accept", "application/json").
		SetResult(&result).
		Get(playtime.featureConfig.APIURL)
	if err != nil {
		return err
	}

	for _, game := range result.Response.Games {
		if game.Playtime == 0 {
			continue
		}
		steamPlaytimeMinutesTotal.WithLabelValues(
			fmt.Sprintf("%d", game.AppID),
			game.Name,
			game.ImgIconURL,
			fmt.Sprintf("http://media.steampowered.com/steamcommunity/public/images/apps/%d/%s.jpg",
				game.AppID, game.ImgIconURL),
			game.ImgLogoURL).
			Set(game.Playtime)
	}

	return nil
}

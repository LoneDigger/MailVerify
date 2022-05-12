package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"me.user/src/config"
	"me.user/src/logger"
	"me.user/src/server"
)

func main() {
	cfg, err := readConfig()
	if err != nil {
		panic(err)
	}

	logger.Init(cfg.Name)

	gin.SetMode(gin.ReleaseMode)

	s := server.NewServer(cfg)
	logrus.WithField("name", cfg.Name).Info("server_start")
	err = s.Start()
	if err != nil {
		logrus.WithError(err).Error()
	}
}

func readConfig() (config.Config, error) {
	var cfg config.Config

	//設定
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yml")

	if err := v.ReadInConfig(); err != nil {
		return cfg, err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

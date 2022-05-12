package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"me.mail/src/config"
	"me.mail/src/logger"
	"me.mail/src/server"
)

func main() {
	cfg, err := readConfig()
	if err != nil {
		panic(err)
	}

	logger.Init(cfg.Name)

	b := make(chan bool)
	s := server.NewServer(cfg)
	s.Start()
	logrus.WithField("name", cfg.Name).Info("server_start")
	<-b
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

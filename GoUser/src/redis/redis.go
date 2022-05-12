package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"me.user/src/config"
)

func newRedis(cfg config.Redis) *redis.Client {
	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Index,
	})

	err := r.Ping().Err()
	if err != nil {
		logrus.WithError(err).Error()
		panic(err)
	}

	return r
}

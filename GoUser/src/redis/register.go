package redis

import (
	"time"

	"github.com/go-redis/redis"
	"me.user/src/config"
)

const registerTime = time.Minute * 3

type Register struct {
	c *redis.Client
}

func NewRegister(cfg config.Redis) *Register {
	return &Register{
		c: newRedis(cfg),
	}
}

func (r *Register) Add(token string) error {
	return r.c.Set(token, "", registerTime).Err()
}

func (r *Register) Check(token string) (bool, error) {
	err := r.c.Get(token).Err()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Register) Remove(token string) {
	r.c.Del(token)
}

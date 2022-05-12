package redis

import (
	"time"

	"github.com/go-redis/redis"
	"me.user/src/config"
)

const validationTime = time.Minute * 30

type Validation struct {
	c *redis.Client
}

func NewValidation(cfg config.Redis) *Validation {
	return &Validation{
		c: newRedis(cfg),
	}
}

func (r *Validation) Add(id int, token string) error {
	return r.c.Set(token, id, validationTime).Err()
}

func (r *Validation) Check(token string) (int, error) {
	return r.c.Get(token).Int()
}

func (r *Validation) Remove(token string) {
	r.c.Del(token)
}

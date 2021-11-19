package helper

import (
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

type Redis struct {
	Url         string
	MaxIdle     int
	IdleTimeout time.Duration
	MaxActive   int
	Pool        *redis.Pool
}

func RegisterRedis() Redis {
	r := Redis{
		Url:         os.Getenv("REDIS"),
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		MaxActive:   10000,
	}
	r.Pool = &redis.Pool{
		MaxIdle:     r.MaxIdle,
		IdleTimeout: r.IdleTimeout,
		MaxActive:   r.MaxActive,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", r.Url)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return r
}

func (r *Redis) Set(key, value string) error {
	_, err := r.Pool.Get().Do("Set", key, value)
	return err
}

func (r *Redis) SetEx(key string, seconds int, value string) error {
	_, err := r.Pool.Get().Do("SetEX", key, seconds, value)
	return err
}

func (r *Redis) Get(key string) (interface{}, error) {
	return r.Pool.Get().Do("Get", key)
}

func (r *Redis) Del(key string) (interface{}, error) {
	return r.Pool.Get().Do("Del", key)
}

func (r *Redis) HSet(name, key string, value interface{}) error {
	_, err := r.Pool.Get().Do("HSet", name, key, value)
	return err
}

func (r *Redis) HGetAll(name string) (map[string]string, error) {
	return redis.StringMap(r.Pool.Get().Do("HGetAll", name))
}

func (r *Redis) HGet(name, key string) (interface{}, error) {
	return r.Pool.Get().Do("HGet", name, key)
}

func (r *Redis) HDel(name, key string) (interface{}, error) {
	return r.Pool.Get().Do("HDel", name, key)
}

func (r *Redis) LPush(name, key string) error {
	_, err := r.Pool.Get().Do("LPush", name, key)
	return err
}

func (r *Redis) RPop(key string) (interface{}, error) {
	return r.Pool.Get().Do("RPop", key)
}

func (r *Redis) GetListLen(key string) (interface{}, error) {
	return r.Pool.Get().Do("LLen", key)
}

func (r *Redis) LRange(key string, start, end int) (interface{}, error) {
	return redis.Values(r.Pool.Get().Do("LRange", key, start, end))
}

func (r *Redis) Expire(key string, seconds int) error {
	_, err := r.Pool.Get().Do("EXPIRE", key, seconds)
	return err
}

package cache

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type redis struct {
	pool *redigo.Pool
}

func (r *redis) Get(k string) (interface{}, bool) {
	conn := r.pool.Get()
	defer conn.Close()

	data, err := redigo.String(conn.Do("GET", k))
	if err != nil {
		return nil, false
	}

	return data, true
}

func (r *redis) Set(k string, v interface{}, d time.Duration) error {
	conn := r.pool.Get()
	defer conn.Close()

	var err error
	if d > DefaultExpiration {
		_, err = conn.Do("SETEX", k, int(d), v)
	} else {
		_, err = conn.Do("SET", k, v)
	}
	if err != nil {
		return errors.WithMessagef(err, "error setting key in redis %s to %s", k, v)
	}

	return nil
}

var onceRedis sync.Once

func (r *redis) initiateRedis() {
	onceRedis.Do(func() {
		redisHost := os.Getenv("REDIS_HOST")
		if redisHost == "" {
			redisHost = ":6379"
		}
		r.newPool(redisHost)
		r.cleanupHook()
	})
}

var redisRef redis

func LoadRedis() Cache {
	if redisRef.pool == nil {
		redisRef.initiateRedis()
	}

	return &redisRef
}

func (r *redis) newPool(server string) {
	r.pool = &redigo.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (r *redis) cleanupHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		err := r.pool.Close()
		if err != nil {
			log.Fatalf("error in cleanupHook(). err: %v", err)
		}
		os.Exit(0)
	}()
}

func (r *redis) ping() error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := redigo.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}

	return nil
}

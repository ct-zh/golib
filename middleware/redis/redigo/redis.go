package redigo

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool
var lock *redsync.Redsync

func Init(host, port, password string) error {
	pool = &redis.Pool{
		MaxIdle:     20,
		IdleTimeout: 240 * time.Second,
		MaxActive:   50,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	lock = redsync.New([]redsync.Pool{pool})
	return nil
}

func GetConn() (redis.Conn, error) {
	if pool == nil {
		return nil, fmt.Errorf("redis is not init")
	}
	conn := pool.Get()
	return conn, conn.Err()
}

func GetLock(key string, expire time.Duration) *redsync.Mutex {
	return lock.NewMutex(key, redsync.SetExpiry(expire))
}

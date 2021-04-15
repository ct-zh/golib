package lock

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/ct-zh/golib/middleware/redis/redigo"
	"github.com/gomodule/redigo/redis"
)

// 分布式锁的简易实现, 解决超卖问题

const fooKey = "test:foo1"

const max = 1000

// 使用redigo附带的分布式锁
func getOne(conn redis.Conn) (bool, error) {
	lock := redigo.GetLock("test:lock", time.Duration(20)*time.Second)
	lock.Lock()
	defer lock.Unlock()
	reply, err := conn.Do("GET", fooKey)
	if err != nil {
		return false, err
	}

	val := 0
	if v, ok := reply.([]byte); ok {
		val, _ = strconv.Atoi(string(v))
		if val >= max {
			return false, errors.New("out of range")
		}
	}

	reply2, err := conn.Do("INCR", fooKey)
	if err != nil {
		return false, err
	}
	log.Printf("get val: %d, reply val: %d", val, reply2)

	return true, nil
}

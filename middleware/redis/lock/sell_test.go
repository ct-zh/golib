package lock

import (
	"github.com/ct-zh/golib/middleware/redis/redigo"
	"github.com/gomodule/redigo/redis"
	"log"
	"sync"
	"testing"
)

//
func TestGetOne(t *testing.T) {
	pool := getPool()
	conn := pool.Get().(redis.Conn)
	conn.Do("DEL", fooKey)

	const WorkerNum = 10

	dataCh := make(chan int, 10)
	wg := sync.WaitGroup{}
	wg.Add(WorkerNum)

	// 生产者
	go func() {
		for i := 0; i < 2000; i++ {
			dataCh <- i
		}
		close(dataCh)
	}()

	// 消费者
	for i := 0; i < WorkerNum; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case _, open := <-dataCh:
					if !open {
						return
					}
					conn := pool.Get().(redis.Conn)
					ok, err := getOne(conn)
					if err != nil {
						t.Fatal("getOne error: ", err)
					}
					if !ok {
						t.Fatal("get one failed")
					}
					pool.Put(conn)
				}
			}
		}()
	}
	wg.Wait()
}

func getPool() *sync.Pool {
	return &sync.Pool{New: func() interface{} {
		redigo.Init("127.0.0.1", "6379", "")
		conn, err := redigo.GetConn()
		if err != nil {
			log.Fatal(err)
		}
		return conn
	}}
}

func doFunc(pool *sync.Pool,
	b *testing.B,
	fn func(redis.Conn) (bool, error)) {

	conn := pool.Get().(redis.Conn)
	ok, err := fn(conn)
	if err != nil {
		b.Fatal("getOne error: ", err)
	}
	if !ok {
		b.Fatal("get one failed")
	}
	pool.Put(conn)
}

func BenchmarkGetOne(b *testing.B) {
	pool := getPool()
	conn := pool.Get().(redis.Conn)
	conn.Do("DEL", fooKey)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doFunc(pool, b, getOne)
		}
	})
}

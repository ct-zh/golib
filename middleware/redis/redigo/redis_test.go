package redigo

import (
	"testing"
	"time"
)

func TestGetConn(t *testing.T) {
	Init("127.0.0.1", "6379", "")
	redis, err := GetConn()
	if err != nil {
		t.Fatal(err)
	}
	defer redis.Close()
	reply, err := redis.Do("SET", "test", 1)
	if err != nil {
		t.Fatal(err)
	}
	switch reply.(type) {
	case []interface{}:
		for i := range reply.([]interface{}) {
			t.Logf("%s", reply.([]interface{})[i])
		}
	default:
		t.Logf("default %+v %T", reply, reply)
	}
}

// BenchmarkGetLock-8 9279 242457 ns/op 844 B/op 24 allocs/op

func BenchmarkGetLock(b *testing.B) {
	err := Init("127.0.0.1", "6379", "")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lockFoo(b)
	}
}

// BenchmarkGetLockParallel-8  100 35361789 ns/op 1907 B/op 31 allocs/op

func BenchmarkGetLockParallel(b *testing.B) {
	err := Init("127.0.0.1", "6379", "")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lockFoo(b)
		}
	})
}

func lockFoo(b *testing.B) {
	lock := GetLock("aaa", 5*time.Second)
	lock.Lock()
	defer lock.Unlock()
	redis, err := GetConn()
	if err != nil {
		b.Error(err)
		return
	}
	defer redis.Close()

	_, err = redis.Do("INCR", "test")
	if err != nil {
		b.Error(err)
		return
	}
}

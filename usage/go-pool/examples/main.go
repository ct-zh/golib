package main

import (
	"fmt"
	"gpool"
	"time"
)

func main() {
	pool, err := gpool.NewPool(3)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 20; i++ {
		pool.Put(&gpool.Task{
			Handler: func(v ...interface{}) {
				fmt.Println(v)
			},
			Params: []interface{}{i},
		})
	}
	time.Sleep(1e5)
}

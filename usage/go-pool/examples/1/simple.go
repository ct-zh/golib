package main

import (
	"fmt"
	"gpool"
	"sync"
)

func main() {
	pool, err := gpool.NewPool(10)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	for i := 0; i < 2000; i++ {
		wg.Add(1)
		task := &gpool.Task{
			Handler: func(v ...interface{}) {
				wg.Done()
				fmt.Println(v)
			},
			Params: []interface{}{i, i * 2, "hello"},
		}
		pool.Put(task)
	}

	wg.Add(1)
	pool.Put(&gpool.Task{
		Handler: func(v ...interface{}) {
			wg.Done()
			fmt.Println(v)
		},
		Params: []interface{}{"bye"},
	})

	wg.Wait()

	pool.Close()
	err = pool.Put(&gpool.Task{
		Handler: func(v ...interface{}) {
			fmt.Println("aaa")
		},
		Params: []interface{}{},
	})
	if err != nil {
		fmt.Println(err)
	}
}

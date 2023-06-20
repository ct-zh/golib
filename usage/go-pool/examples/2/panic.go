package main

import (
	"fmt"
	"gpool"
)

func main() {
	pool, err := gpool.NewPool(10)
	if err != nil {
		panic(err)
	}
	pool.PanicHandler = func(i interface{}) {
		fmt.Println("warning!!! panic")
	}

	pool.Put(&gpool.Task{
		Handler: func(v ...interface{}) {
			panic("eeee")
		},
		Params: nil,
	})
}

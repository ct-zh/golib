package threading

import "log"

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if r := recover(); r != nil {
		log.Print(r)
	}
}

func RunSafe(fn func()) {
	defer Recover()
	fn()
}

func GoSafe(fn func()) {
	go RunSafe(fn)
}

package threading

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if r := recover(); r != nil {
		panic(r) // ?
	}
}

func RunSafe(fn func()) {
	defer Recover()
	fn()
}

func GoSafe(fn func()) {
	go RunSafe(fn)
}

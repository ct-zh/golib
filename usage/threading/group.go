package threading

import "sync"

// A RoutineGroup is used to group goroutines together
// and all wait all goroutines to be done.
type RoutineGroup struct {
	wg sync.WaitGroup
}

func NewRoutineGroup() *RoutineGroup {
	return &RoutineGroup{}
}

// Run runs the gives function in RoutineGroup.
// Please don't reference any variables from outside
// because outside variables can be changed by other goroutines
func (g *RoutineGroup) Run(fn func()) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()
		fn()
	}()
}

func (g *RoutineGroup) RunSafe(fn func()) {
	g.wg.Add(1)
	GoSafe(func() {
		defer g.wg.Done()
		fn()
	})
}

func (g *RoutineGroup) Wait() {
	g.wg.Wait()
}

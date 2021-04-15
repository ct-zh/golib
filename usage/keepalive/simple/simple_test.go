package simple

import "testing"

func TestNewWorkerManager(t *testing.T) {
	work := NewWorkerManager(10)
	work.StartWorkerPool()
}

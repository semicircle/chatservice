package chatservice

import (
	"testing"
)

func newBaseThread() Thread {
	return NewBaseThread()
}

func TestBaseThread(t *testing.T) {
	threadtest1(t, newBaseThread)
}

func BenchmarkBaseThread1(b *testing.B) {
	threadbench1(b, newBaseThread)
}

func BenchmarkBaseThread2(b *testing.B) {
	threadbench2(b, newBaseThread)
}

func BenchmarkBaseThread3(b *testing.B) {
	threadbench3(b, newBaseThread)
}

package chatservice

import (
	"testing"
)

func newBaseMessage() Message {
	return NewBaseMessage()
}

func TestBaseMessage(t *testing.T) {
	messagetest1(t, newBaseMessage)
}

func BenchmarkBaseMesage1(b *testing.B) {
	messagebench1(b, newBaseMessage)
}

func BenchmarkBaseMesage2(b *testing.B) {
	messagebench2(b, newBaseMessage)
}

func BenchmarkBaseMesage3(b *testing.B) {
	messagebench3(b, newBaseMessage)
}

package chatservice

import (
	"testing"
)

func newBaseUser() User {
	return NewBaseUser()
}

func TestBaseUser1(t *testing.T) {
	usertest1(t, newBaseUser)
}

func TestBaseUser2(t *testing.T) {
	usertest2(t, newBaseUser)
}

func TestBaseUser3(t *testing.T) {
	usertest3(t, newBaseUser)
}

func BenchmarkBaseUser1(b *testing.B) {
	userbench1(b, newBaseUser)
}

func BenchmarkBaseUser2(b *testing.B) {
	userbench2(b, newBaseUser)
}

func BenchmarkBaseUser3(b *testing.B) {
	userbench3(b, newBaseUser)
}

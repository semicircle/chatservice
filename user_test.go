package chatservice

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

//common test functions of User interface

type newuserfunc func() User

func usertest1(t *testing.T, newuser newuserfunc) {
	u := newuser()
	u.SetNick("semicircle")
	if u.Nick() != "semicircle" {
		t.Error("SetNick failed")
	}
	u.SetPrivateKey(keyType(0x12345678))
	if u.PrivateKey() != keyType(0x12345678) {
		t.Error("SetPrivateKey failed")
	}
}

func usertest2(t *testing.T, newuser newuserfunc) {
	u1 := newuser()
	u2 := newuser()

	u1.Generate()
	u2.Generate()

	if u1.PrivateKey() == u2.PrivateKey() {
		t.Error("Generate failed: PrivateKey() is not indentical")
	}
}

func randuser(newuser newuserfunc) User {
	ret := newuser()
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret.Generate()
	ret.SetNick(strconv.Itoa(r.Int())).SetPrivateKey(keyType(r.Int63()))
	return ret
}

func usertest3(t *testing.T, newuser newuserfunc) {
	entitytest(t, func() Entity {
		return randuser(newuser)
	}, func(a, b Entity) bool {
		var x, y User
		x, _ = a.(User)
		y, _ = b.(User)
		return (x.Nick() == y.Nick()) && (x.PrivateKey() == y.PrivateKey())
	})
}

func userbench1(b *testing.B, newuser newuserfunc) {
	entitybench1(b, func() Entity {
		b.StopTimer()
		defer b.StartTimer()
		return randuser(newuser)
	})
}

func userbench2(b *testing.B, newuser newuserfunc) {
	entitybench2(b, func() Entity {
		b.StopTimer()
		defer b.StartTimer()
		return randuser(newuser)
	})
}

func userbench3(b *testing.B, newuser newuserfunc) {
	entitybench3(b, func() Entity {
		b.StopTimer()
		defer b.StartTimer()
		return randuser(newuser)
	})
}

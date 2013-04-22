package chatservice

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type newthreadfunc func() Thread

func randthread(newthread newthreadfunc) Thread {
	ret := newthread()
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret.SetTitle(strconv.Itoa(r.Int()))
	return ret
}

func threadtest1(t *testing.T, newthread newthreadfunc) {
	entitytest(t, func() Entity {
		return randthread(newthread)
	}, func(a, b Entity) bool {
		var x, y Thread
		x, _ = a.(Thread)
		y, _ = b.(Thread)
		ret := x.Title() == y.Title()
		return ret
	})
}

func threadbench1(b *testing.B, newthread newthreadfunc) {
	entitybench1(b, func() Entity {
		b.StopTimer()
		defer b.StartTimer()
		return randthread(newthread)
	})
}

func threadbench2(b *testing.B, newthread newthreadfunc) {
	entitybench2(b, func() Entity {
		b.StopTimer()
		defer b.StartTimer()
		return randthread(newthread)
	})
}

func threadbench3(b *testing.B, newthread newthreadfunc) {
	entitybench3(b, func() Entity {
		b.StopTimer()
		defer b.StartTimer()
		return randthread(newthread)
	})
}

package chatservice

import (
	"testing"
)

type newentityfunc func() Entity
type entitycheckequalfunc func(a, b Entity) bool

func entitytest(t *testing.T, newentiy newentityfunc, checkequal entitycheckequalfunc) {
	u := newentiy()
	defer u.Delete()
	if err := u.Save(); err != nil {
		t.Error("Save failed")
	}

	u2 := newentiy()
	defer u2.Delete()
	if err := u2.Save(); err != nil {
		t.Error("Save failed")
	}

	if u.Id() == u2.Id() {
		t.Error("Save failed: id dup")
	}

	u3 := newentiy()
	if err := u3.Load(u.Id()); err != nil {
		t.Error("Load failed")
	}

	if !checkequal(u, u3) {
		t.Error("Save/Load Error: checkequal")
	}

	if err := u3.Delete(); err != nil {
		t.Error("Delete failed")
	}

	u4 := newentiy()
	if err := u4.Load(u.Id()); err == nil {
		t.Error("Delete/Load failed")
	}

}

func entitybench1(b *testing.B, newentiy newentityfunc) {
	b.StopTimer()
	benchmarkusers := make([]Entity, b.N)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u := newentiy()
		u.Save()
		benchmarkusers[i] = u
	}
}

func entitybench2(b *testing.B, newentiy newentityfunc) {
	b.StopTimer()
	benchmarkusers := make([]Entity, b.N)
	for i := 0; i < b.N; i++ {
		u := newentiy()
		u.Save()
		benchmarkusers[i] = u
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u := newentiy()
		u.Load(benchmarkusers[i].Id())
	}
}

func entitybench3(b *testing.B, newentiy newentityfunc) {
	b.StopTimer()
	benchmarkusers := make([]Entity, b.N)
	for i := 0; i < b.N; i++ {
		u := newentiy()
		u.Save()
		benchmarkusers[i] = u
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		benchmarkusers[i].Delete()
	}
}

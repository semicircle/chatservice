package chatservice

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type newmsgfunc func() Message

func randmessage(newmsg newmsgfunc) Message {
	ret := newmsg()
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret.SetAuthorId(idType(r.Int()))
	ret.SetDisplyAuthor(strconv.Itoa(r.Int()))
	ret.SetText(strconv.Itoa(r.Int()))
	ret.SetType(MessageCode(r.Int() % 3))
	ret.SetDate(time.Unix(r.Int63(), r.Int63()))
	return ret
}

func messagetest1(t *testing.T, newmsg newmsgfunc) {
	entitytest(t, func() Entity {
		return randmessage(newmsg)
	}, func(a, b Entity) bool {
		var x, y Message
		x, _ = a.(Message)
		y, _ = b.(Message)
		ret := (x.AuthorId() == y.AuthorId()) && (x.Date() == y.Date()) && (x.DisplyAuthor() == y.DisplyAuthor()) && (x.Text() == y.Text()) && (x.Type() == y.Type())
		return ret
	})
}

func messagebench1(b *testing.B, newmsg newmsgfunc) {
	entitybench1(b, func() Entity {
		b.StopTimer()
		defer b.StartTimer()
		return randmessage(newmsg)
	})
}

func messagebench2(b *testing.B, newmsg newmsgfunc) {
	entitybench2(b, func() Entity {
		b.StopTimer()
		defer b.StartTimer()
		return randmessage(newmsg)
	})
}

func messagebench3(b *testing.B, newmsg newmsgfunc) {
	entitybench3(b, func() Entity {
		b.StopTimer()
		defer b.StartTimer()
		return randmessage(newmsg)
	})
}

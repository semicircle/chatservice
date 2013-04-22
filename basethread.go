package chatservice

import (
	"errors"
)

// type Thread interface {
// 	Entity
// 	Subscribe(u User) error
// 	Unsubscribe(u User) error
// 	Post(u User, m Message) error
// 	Achieve() []Message
// }

const (
	BASETHREAD_BUFFER_SIZE   = 4096
	BASETHREAD_ARCHIEVE_SIZE = 100
	BASETHREAD_SWITCHER_NUM  = 20
)

type BaseThread struct {
	BaseEntity
	subscribers map[idType]User
	posting     chan Message
	archieve    []Message
	release     chan bool
	title       string
}

var (
	basethreadstore = NewBaseStore()
)

func NewBaseThread() *BaseThread {
	ret := &BaseThread{
		NewBaseEntity(basethreadstore, func() Entity {
			return NewBaseThread()
		}, cloneBaseThread),
		make(map[idType]User),
		make(chan Message, BASETHREAD_BUFFER_SIZE),
		make([]Message, 0, BASETHREAD_ARCHIEVE_SIZE),
		make(chan bool, BASETHREAD_SWITCHER_NUM*10),
		""}
	for i := 0; i < BASETHREAD_SWITCHER_NUM; i++ {
		go ret.messageswitcher()
	}
	return ret
}

func (t *BaseThread) Title() string {
	return t.title
}

func (t *BaseThread) SetTitle(title string) Thread {
	t.title = title
	return t
}

func (t *BaseThread) Subscribe(u User) error {
	if _, ok := t.subscribers[u.Id()]; !ok {
		t.subscribers[u.Id()] = u
	} else {
		return errors.New("Already subscribed.")
	}
	return nil
}

func (t *BaseThread) Unsubscribe(u User) error {
	if _, ok := t.subscribers[u.Id()]; !ok {
		return errors.New("Not subscribed.")
	} else {
		delete(t.subscribers, u.Id())
	}
	return nil
}

func (t *BaseThread) Post(u User, m Message) error {
	t.archieve = append(t.archieve, m)
	t.posting <- m
	return nil
}

func (t *BaseThread) Achieve() []Message {
	return t.archieve
}

func (t *BaseThread) messageswitcher() {
	for {
		select {
		case msg := <-t.posting:
			for _, u := range t.subscribers {
				// go func(u User, msg Message) {
				// 	for times := 0; times < 10; times++ {
				// 		if err := u.WakeRecvPending(msg); err != nil {
				// 			if times < 120 {
				// 				time.Sleep(time.Duration(int(time.Second) * times))
				// 			} else {
				// 				break //too many times
				// 			}
				// 		} else {
				// 			break //retry success
				// 		}
				// 	}
				// }(u, msg)
				u.WakeRecvPending(msg)
			}
		case <-t.release:
			return
		}
	}
}

func (t *BaseThread) Release() {
	for i := 0; i < BASETHREAD_SWITCHER_NUM*2; i++ {
		if len(t.release) < BASETHREAD_SWITCHER_NUM*10 {
			t.release <- true
		}
	}
}

//Storage Part.
func cloneBaseThread(dst, src Entity) {
	var (
		d, s *BaseThread
		ok   bool
	)
	if d, ok = dst.(*BaseThread); !ok {
		panic("cloneBaseThread called with wrong dst Entity.")
	}
	if s, ok = src.(*BaseThread); !ok {
		panic("cloneBaseThread called with wrong src Entity.")
	}
	d.subscribers = s.subscribers
	d.posting = s.posting
	d.archieve = s.archieve
	d.release = s.release
}

func (t *BaseThread) Save() error {
	return t.SaveEntity(t)
}

func (t *BaseThread) Load(id idType) error {
	return t.LoadEntity(t, id)
}

func (t *BaseThread) Delete() error {
	t.Release()
	return t.BaseEntity.Delete()
}

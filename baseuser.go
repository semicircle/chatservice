package chatservice

import (
	"errors"
	"sync"
	"time"
)

const (
	RECV_BUF = 1000
)

type BaseUser struct {
	BaseEntity
	nick       string
	privateKey keyType
	recvbuf    chan Message
}

func NewBaseUser() *BaseUser {
	return &BaseUser{
		BaseEntity{idType(0), baseuserstore, func() Entity {
			return NewBaseUser()
		}, cloneBaseUser},
		"", keyType(0), make(chan Message, RECV_BUF)}

}

// Attributes
func (u *BaseUser) Nick() string {
	return u.nick
}

func (u *BaseUser) PrivateKey() keyType {
	return u.privateKey

}

func (u *BaseUser) SetNick(nick string) User {
	u.nick = nick
	return u
}

func (u *BaseUser) SetPrivateKey(key keyType) User {
	u.privateKey = key
	return u
}

// Realtime Part.
func (u *BaseUser) WakeRecvPending(msg Message) error {
	if len(u.recvbuf) < RECV_BUF {
		u.recvbuf <- msg
		return nil
	}
	return errors.New("BufferFull")
}

func (u *BaseUser) WaitRecvMessage(timeout <-chan time.Time) []Message {
	var msg Message
	ret := make([]Message, 0, 10)
	select {
	case <-timeout:
		return nil
	case msg = <-u.recvbuf:
		ret = append(ret, msg)
		for len(u.recvbuf) != 0 {
			msg = <-u.recvbuf
			ret = append(ret, msg)
		}
	}

	return ret
}

// Storage Part.
type keysInfo struct {
	keyindex keyType
	mutex    sync.Mutex
}

var (
	keysinfo      = keysInfo{keyType(0), sync.Mutex{}}
	baseuserstore = NewBaseStore()
)

func (u *BaseUser) Generate() User {
	keysinfo.mutex.Lock()
	defer keysinfo.mutex.Unlock()
	u.privateKey = keysinfo.keyindex
	keysinfo.keyindex++

	return u
}

func cloneBaseUser(dst, src Entity) {
	var (
		d, s *BaseUser
		ok   bool
	)
	if d, ok = dst.(*BaseUser); !ok {
		panic("cloneBaseUser called with wrong dst Entity.")
	}
	if s, ok = src.(*BaseUser); !ok {
		panic("cloneBaseUser called with wrong src Entity.")
	}
	d.nick = s.nick
	d.privateKey = s.privateKey
	d.recvbuf = s.recvbuf
}

func (u *BaseUser) Save() error {
	return u.SaveEntity(u)
}

func (u *BaseUser) Load(id idType) error {
	return u.LoadEntity(u, id)
}

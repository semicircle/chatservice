package chatservice

import (
	"time"
)

type User interface {
	Entity
	Nick() string
	PrivateKey() keyType
	Generate() User //Generate User Id / PrivateKey
	SetNick(nick string) User
	SetPrivateKey(key keyType) User

	WakeRecvPending(msg Message) error
	WaitRecvMessage(timeout <-chan time.Time) []Message
}

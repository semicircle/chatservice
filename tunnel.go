package chatservice

import (
	"time"
)

type Tunnel interface {
	RecvFromClient(timeout chan time.Time) (User, []byte, error)
	SendToClient(u User, jsonblob []byte, timeout <-chan time.Time) error
}

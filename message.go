package chatservice

import (
	"time"
)

type MessageCode int

const (
	MESSAGE_TYPE_NORMAL      = iota // thread messages
	MESSAGE_TYPE_NOTIFY             // join / leave thread.
	MESSAGE_TYPE_SYSTEM             // "the system is going to reboot"
	MESSAGE_TYPE_USERPRIVATE        // "hi~"
)

type Message interface {
	Entity
	Type() MessageCode
	SetType(t MessageCode) Message
	Text() string
	SetText(text string) Message
	DisplyAuthor() string
	SetDisplyAuthor(text string) Message
	AuthorId() idType
	SetAuthorId(id idType) Message
	Date() time.Time
	SetDate(time time.Time) Message
	ThreadId() idType
	SetThreadId(id idType) Message
}

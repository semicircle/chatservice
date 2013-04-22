package chatservice

import (
	"time"
)

type MessageCode int

const (
	MESSAGE_TYPE_NORMAL      = 0 // thread messages
	MESSAGE_TYPE_NOTIFY      = 1 // join / leave thread.
	MESSAGE_TYPE_SYSTEM      = 2 // "the system is going to reboot"
	MESSAGE_TYPE_USERPRIVATE = 3 // "hi~"
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

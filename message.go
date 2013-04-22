package chatservice

import (
	"time"
)

type MessageCode int

const (
	MESSAGE_TYPE_NORMAL = 0
	MESSAGE_TYPE_NOTIFY = 1
	MESSAGE_TYPE_SYSTEM = 2
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
}

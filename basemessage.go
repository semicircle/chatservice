package chatservice

import (
	"time"
)

type BaseMessage struct {
	BaseEntity
	msgCode       MessageCode
	text          string
	displayAuthor string
	authorId      idType
	date          time.Time
}

func NewBaseMessage() *BaseMessage {
	return &BaseMessage{
		NewBaseEntity(basemessagestore, func() Entity {
			return NewBaseMessage()
		}, cloneBaseMessage),
		MessageCode(0), "", "", idType(0), time.Time{}}
}

//Attributes.
func (m *BaseMessage) Type() MessageCode {
	return m.msgCode
}

func (m *BaseMessage) SetType(c MessageCode) Message {
	m.msgCode = c
	return m
}

func (m *BaseMessage) Text() string {
	return m.text
}

func (m *BaseMessage) SetText(text string) Message {
	m.text = text
	return m
}

func (m *BaseMessage) DisplyAuthor() string {
	return m.displayAuthor
}

func (m *BaseMessage) SetDisplyAuthor(displyAuthor string) Message {
	m.displayAuthor = displyAuthor
	return m
}

func (m *BaseMessage) AuthorId() idType {
	return m.authorId
}

func (m *BaseMessage) SetAuthorId(authorId idType) Message {
	m.authorId = authorId
	return m
}

func (m *BaseMessage) Date() time.Time {
	return m.date
}

func (m *BaseMessage) SetDate(date time.Time) Message {
	m.date = date
	return m
}

//Storage Part.
var (
	basemessagestore = NewBaseStore()
)

func cloneBaseMessage(dst, src Entity) {
	var (
		d, s *BaseMessage
		ok   bool
	)
	if d, ok = dst.(*BaseMessage); !ok {
		panic("cloneBaseUser called with wrong dst Entity.")
	}
	if s, ok = src.(*BaseMessage); !ok {
		panic("cloneBaseUser called with wrong src Entity.")
	}
	d.msgCode = s.msgCode
	d.text = s.text
	d.displayAuthor = s.displayAuthor
	d.authorId = s.authorId
	d.date = s.date
}

func (m *BaseMessage) Save() error {
	return m.SaveEntity(m)
}

func (m *BaseMessage) Load(id idType) error {
	return m.LoadEntity(m, id)
}

package chatservice

type Thread interface {
	Entity
	Title() string
	SetTitle(title string) Thread
	Subscribe(u User) error
	Unsubscribe(u User) error
	Post(u User, m Message) error
	Achieve() []Message
	Release()
}

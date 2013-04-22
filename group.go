package chatservice

type Group interface {
	Entity
	Threads() []Thread
	Users() []User
	Join(u User)
	Leave(u User)
}

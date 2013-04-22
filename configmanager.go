package chatservice

type Factory struct {
	User    func() User
	Thread  func() Thread
	Message func() Message
	Group   func() Group
}

var (
	FACTORY = Factory{
		func() User {
			return NewBaseUser()
		},
		func() Thread {
			return NewBaseThread()
		},
		func() Message {
			return NewBaseMessage()
		}, nil}
)

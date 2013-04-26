package chatservice

import (
	"fmt"
	"time"
)

//ChatService is a http.Handler
type ChatService struct {
	sigProcessor *SigProcessor
	tunnel       Tunnel
}

func NewChatService(t Tunnel) (cs *ChatService) {
	cs = &ChatService{NewSigProcessor(), t}
	go cs.chatServeRountine()
	go cs.chatServeRountine()
	go cs.chatServeRountine()
	return
}

func (cs *ChatService) chatServeRountine() {
	for {
		u, jsonblob, err := cs.tunnel.RecvFromClient(make(chan time.Time))
		if err != nil {
			fmt.Println(err)
			return
		}
		retblob, err2 := cs.sigProcessor.Handle(u, jsonblob)
		if err2 != nil {
			fmt.Println(err)
			return
		}
		if err3 := cs.tunnel.SendToClient(u, retblob, time.After(time.Minute)); err3 != nil {
			fmt.Println(err)
			return
		}
	}
}

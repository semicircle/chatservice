package main

import (
	"fmt"
	"github.com/semicircle/chatservice"
)

func main() {
	fmt.Println("halo")
	chatservice.NewChatService(chatservice.NewSimpleTunnel(""))
	return
}

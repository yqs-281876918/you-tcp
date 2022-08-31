package main

import (
	"fmt"
	"youtcp/you_interface"
	"youtcp/you_net"
)

type PingReqHandler struct {
	*you_net.DefaultRequestHandler
}

func (handler *PingReqHandler) Handle(request you_interface.IYouRequest) {
	fmt.Printf("receive message : %v\n", request.GetMessage().GetData())
}

func main() {
	server := you_net.NewServer()
	server.RegisterHandler(1, &PingReqHandler{&you_net.DefaultRequestHandler{}})
	server.Run()
}

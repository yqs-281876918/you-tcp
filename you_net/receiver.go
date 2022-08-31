package you_net

type RequestReceiver struct {
	toReceive    int
	handlerMap   map[uint32]RequestHandler
	reqChan      []chan *YouRequest
	workPoolSize int
	started      bool
}

func (this *RequestReceiver) RegisterHandler(msgID uint32, handler RequestHandler) {
	this.handlerMap[msgID] = handler
}

func (this *RequestReceiver) getHandler(msgID uint32) RequestHandler {
	return this.handlerMap[msgID]
}

func (this *RequestReceiver) AddRequest(request *YouRequest) {
	this.reqChan[this.toReceive] <- request
	this.toReceive = (this.toReceive + 1) % this.workPoolSize
}

func NewReqHandlerReceiver() *RequestReceiver {
	rhr := &RequestReceiver{
		toReceive:    0,
		handlerMap:   make(map[uint32]RequestHandler),
		workPoolSize: 16,
	}
	rhr.reqChan = make([]chan *YouRequest, rhr.workPoolSize)
	for i := 0; i < rhr.workPoolSize; i++ {
		rhr.reqChan[i] = make(chan *YouRequest, 1024*16)
	}
	rhr.start()
	return rhr
}

func (this *RequestReceiver) start() {
	if this.started {
		return
	}
	this.started = true
	for i := 0; i < this.workPoolSize; i++ {
		go this.handleRequest(i)
	}
}

func (this *RequestReceiver) handleRequest(num int) {
	for {
		request := <-this.reqChan[num]
		handler := this.getHandler(request.GetMessage().GetID())
		if handler == nil {
			continue
		}
		handler.PreHandle(request)
		handler.Handle(request)
		handler.PostHandle(request)
	}
}

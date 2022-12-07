package you_net

import (
	"fmt"
	"net"
	"youtcp/config"
)

type Server struct {
	name              string //服务器名称
	ipVersion         string //ip版本
	ip                string //监听ip
	port              int    //监听端口
	nextConnID        uint32
	requestReceiver   *RequestReceiver
	connManager       *ConnManager
	onConnStartedFunc func(conn *YouConnection)
	onConnStopFunc    func(conn *YouConnection)
}

func (server *Server) start() bool {
	fmt.Printf("server start at ip:%s port:%d\n", server.ip, server.port)
	addr, err := net.ResolveTCPAddr(server.ipVersion, fmt.Sprintf("%s:%d", server.ip, server.port))
	if err != nil {
		fmt.Printf("resovle tcp address error : %v\n", err)
		return false
	}
	//获取监听对象
	listener, err := net.ListenTCP(server.ipVersion, addr)
	if err != nil {
		fmt.Printf("listen error : %v\n", err)
		return false
	}
	fmt.Println("server start successfully")
	//开始监听客户端连接
	go server.listen(listener)
	return true
}

func (server *Server) listen(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("accept connection error : %v\n", err)
			continue
		}
		server.handleConn(conn)
	}
}

func (server *Server) handleConn(conn *net.TCPConn) {
	youConn := NewConnection(conn, server.nextConnID, server)
	server.nextConnID++
	server.connManager.AddConn(youConn)
}

func (server *Server) Stop() {

}

func (server *Server) Run() {
	if !server.start() {
		fmt.Println("server start fail")
		return
	}
	//阻塞住
	select {}
}

func (server *Server) RegisterHandler(id uint32, requestHandler RequestHandler) {
	server.requestReceiver.RegisterHandler(id, requestHandler)
}

func (server *Server) GetRequestReceiver() *RequestReceiver {
	return server.requestReceiver
}

func (server *Server) GetConnManager() *ConnManager {
	return server.connManager
}

func (server *Server) OnConnStarted(fun func(conn *YouConnection)) {
	server.onConnStartedFunc = fun
}

func (server *Server) OnConnStop(fun func(conn *YouConnection)) {
	server.onConnStopFunc = fun
}

func NewServer() *Server {
	cfg := config.GetGlobalConfig()
	server := &Server{
		name:            cfg.Name,
		ipVersion:       "tcp4",
		ip:              cfg.Host,
		port:            cfg.Port,
		nextConnID:      0,
		requestReceiver: NewReqHandlerReceiver(),
		connManager:     NewConnManager(),
	}
	return server
}

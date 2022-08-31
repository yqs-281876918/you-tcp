package you_net

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
)

type YouConnection struct {
	conn       *net.TCPConn
	connID     uint32
	writerChan chan *YouMessage
	server     *Server
	properties map[string]interface{}
	rwMutex    sync.RWMutex
}

func (conn *YouConnection) startReader() {
	defer conn.Close()
	for {
		msg := new(YouMessage)
		headBuffer := make([]byte, 8)
		_, err := io.ReadFull(conn.conn, headBuffer)
		if err != nil {
			return
		}
		msgLen := binary.BigEndian.Uint32(headBuffer[0:4])
		msgID := binary.BigEndian.Uint32(headBuffer[4:8])
		msg.id = msgID
		contentBuffer := make([]byte, msgLen)
		_, err = io.ReadFull(conn.conn, contentBuffer)
		if err != nil {
			return
		}
		msg.data = contentBuffer
		if err != nil {
			fmt.Printf("reader error = %v\n", err)
			return
		}
		request := &YouRequest{
			conn: conn,
			msg:  msg,
		}
		conn.server.GetRequestReceiver().AddRequest(request)
	}
}

func (conn *YouConnection) startWriter() {
	defer conn.Close()
	for {
		msg, ok := <-conn.writerChan
		if !ok {
			return
		}
		lenBytes := make([]byte, 4)
		idBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(lenBytes, msg.GetLen())
		binary.BigEndian.PutUint32(idBytes, msg.GetID())
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("消息序列化失败,暂停发送")
			return
		}
		if !conn.writeBytes(lenBytes) || !conn.writeBytes(idBytes) || !conn.writeBytes(msgBytes) {
			fmt.Printf("消息发送失败,暂停发送")
			return
		}
	}
}

func (conn *YouConnection) Start() {
	fmt.Printf("conn start connID=%d\n", conn.connID)
	//读数据
	go conn.startReader()
	//写数据
	go conn.startWriter()
	if conn.server.onConnStopFunc != nil {
		conn.server.onConnStartedFunc(conn)
	}
}

func (conn *YouConnection) Close() {
	conn.server.GetConnManager().RemoveConn(conn)
	err := conn.conn.Close()
	if err != nil {
		return
	}
	if conn.server.onConnStopFunc != nil {
		conn.server.onConnStopFunc(conn)
	}
}

func (conn *YouConnection) GetConnectionID() uint32 {
	return conn.connID
}
func (conn *YouConnection) RemoteAddress() net.Addr {
	return conn.conn.RemoteAddr()
}
func (conn *YouConnection) SendMessage(msg *YouMessage) {
	conn.writerChan <- msg
}

func (conn *YouConnection) SetProperties(key string, value interface{}) {
	conn.rwMutex.Lock()
	defer conn.rwMutex.Unlock()
	conn.properties[key] = value
}

func (conn *YouConnection) GetProperties(key string) interface{} {
	conn.rwMutex.RLock()
	defer conn.rwMutex.RLock()
	if v, ok := conn.properties[key]; ok {
		return v
	} else {
		return nil
	}
}

func (conn *YouConnection) RemoveProperties(key string) {
	conn.SetProperties(key, nil)
}

func (conn *YouConnection) writeBytes(data []byte) bool {
	nbFailed := 0
	for {
		if nbFailed == 10 {
			return false
		}
		n, err := conn.conn.Write(data)
		if err != nil {
			nbFailed++
			fmt.Println("conn发送数据失败,尝试重发...")
			data = data[n:]
			continue
		}
		break
	}
	return true
}

func NewConnection(conn *net.TCPConn, connID uint32, server *Server) *YouConnection {
	return &YouConnection{
		conn:       conn,
		connID:     connID,
		writerChan: make(chan *YouMessage, 16),
		server:     server,
		properties: make(map[string]interface{}),
	}
}

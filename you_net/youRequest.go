package you_net

type YouRequest struct {
	conn *YouConnection
	msg  *YouMessage
}

func (request *YouRequest) GetConnection() *YouConnection {
	return request.conn

}

func (request *YouRequest) GetMessage() *YouMessage {
	return request.msg
}

package you_net

type YouMessage struct {
	id   uint32
	data []byte
}

func (msg *YouMessage) GetID() uint32 {
	return msg.id
}
func (msg *YouMessage) GetData() []byte {
	return msg.data
}
func (msg *YouMessage) GetLen() uint32 {
	return uint32(len(msg.data))
}

func NewYouMessage(id uint32, data []byte) *YouMessage {
	return &YouMessage{id: id, data: data}
}

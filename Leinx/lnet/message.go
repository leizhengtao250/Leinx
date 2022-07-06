package lnet

type Message struct {
	Id      uint32 //消息ID
	Datalen uint32 //消息长度
	Data    []byte //消息内容
}

func (m *Message) GetMsgID() uint32 {
	return m.Id
}

func (m *Message) GetMsgLen() uint32 {
	return m.Datalen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsID(id uint32) {
	m.Id = id
}

func (m *Message) SetMsgLen(len uint32) {
	m.Datalen = len
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

//创建一个Message消息包

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		Datalen: uint32(len(data)),
		Data:    data,
	}

}

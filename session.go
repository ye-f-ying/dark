package dark

type Session struct {
	c       *Conn
	msg     []byte
	msgType int
}

func NewSession(c *Conn, msg []byte, msgType int) *Session {
	return &Session{c: c, msg: msg, msgType: msgType}
}

/**
 * @description: 获取链接
 * @return {*}
 */
func (m *Session) GetConn() *Conn {
	return m.c
}

/**
 * @description: 获取消息
 * @return {*}
 */
func (m *Session) GetMsg() []byte {
	return m.msg
}

/**
 * @description: 消息类型--websocket 才有
 * @return {*}
 */
func (m *Session) GetMsgType() int {
	return m.msgType
}

/**
 * @description: 发送消息
 * @param {[]byte} data
 * @return {*}
 */
func (m *Session) Send(data []byte) error {
	if m.c == nil {
		return nil
	}
	if m.c.GetService().GetOptions().IsWebsocket {
		return m.c.GetWSConn().WriteMessage(m.msgType, data)
	}
	_, err := m.c.GetConn().Write(data)
	return err

}

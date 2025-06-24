/*
 * @Author: yeying
 * @Date: 2025-06-05 22:27:56
 * @LastEditTime: 2025-06-23 22:18:03
 * @FilePath: \dark\session.go
 * @Description:
 */
package dark

import "context"

type Session struct {
	c    *Conn
	pack IPack
	ctx  context.Context
}

func NewSession(c *Conn, p IPack, ctx context.Context) *Session {
	return &Session{c: c, pack: p, ctx: ctx}
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
	if m.pack == nil {
		return nil
	}
	return m.pack.GetData()
}

/**
 * @description: 消息类型--websocket 才有
 * @return {*}
 */
func (m *Session) GetMsgType() int {
	if m.pack == nil {
		return 0
	}
	return m.pack.GetType()
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
		return m.c.GetWSConn().WriteMessage(m.pack.GetType(), data)
	}
	_, err := m.c.GetConn().Write(data)
	return err

}

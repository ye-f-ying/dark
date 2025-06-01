/*
 * @Author: yeying
 * @Date: 2025-06-01 18:05:29
 * @LastEditTime: 2025-06-01 18:08:47
 * @FilePath: \dark\conn.go
 * @Description:
 */
package dark

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/panjf2000/gnet/v2"
)

type Conn struct {
	ctx    context.Context
	wsConn *websocket.Conn //websocket链接
	c      gnet.Conn
}

/**
 * @description: 创建一个链接
 * @param {gnet.Conn} c
 * @return {*}
 */
func NewConn(c gnet.Conn) *Conn {
	return &Conn{ctx: context.Background(), wsConn: nil, c: c}
}

/**
 * @description: 获取context
 * @return {*}
 */
func (m *Conn) GetContext() context.Context {
	return m.ctx
}

/**
 * @description: 设置websocket链接
 * @param {*websocket.Conn} wsConn
 * @return {*}
 */
func (m *Conn) SetWSConn(wsConn *websocket.Conn) {
	m.wsConn = wsConn
}

/**
 * @description: 设置GNET链接
 * @param {gnet.Conn} c
 * @return {*}
 */
func (m *Conn) SetConn(c gnet.Conn) {
	m.c = c
}

/**
 * @description: 获取websocket链接
 * @param {*websocket.Conn} wsConn
 * @return {*}
 */
func (m *Conn) GetWSConn() *websocket.Conn {
	return m.wsConn
}

/**
 * @description:读取数据
 * @return {*}
 */
func (m *Conn) Read() (int, []byte, error) {
	if m.c == nil {
		return 0, nil, fmt.Errorf("error conn")
	}
	if m.wsConn != nil {
		return m.wsConn.ReadMessage()
	}

	size := m.c.InboundBuffered()
	buf := make([]byte, size)
	read, err := m.c.Read(buf)
	if err != nil || read < size {
		return read, nil, err
	}
	return read, buf, nil
}

/**
 * @description:发送数据
 * @param {[]byte} buf
 * @param {...interface{}} argv
 * @return {*}
 */
func (m *Conn) Send(buf []byte, msgTypes ...int) (int, error) {
	if m.c == nil {
		return -1, fmt.Errorf("error conn")
	}
	bufLen := len(buf)
	if bufLen <= 0 {
		return -1, fmt.Errorf("data is null")
	}

	if m.wsConn != nil {
		messageType := websocket.BinaryMessage
		if len(msgTypes) > 0 {
			messageType = msgTypes[0]
		}
		return bufLen, m.wsConn.WriteMessage(messageType, buf)
	}

	return m.c.Write(buf)
}

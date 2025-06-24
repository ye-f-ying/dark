/*
 * @Author: yeying
 * @Date: 2025-06-01 18:05:29
 * @LastEditTime: 2025-06-24 22:38:15
 * @FilePath: \dark\conn.go
 * @Description:
 */
package dark

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/panjf2000/gnet/v2"
)

type Conn struct {
	ctx        context.Context
	wsConn     *websocket.Conn //websocket链接
	c          gnet.Conn
	activeTime int64 //最近活跃时间
	service    *Service
	connTime   int64       //连接时间
	mark       interface{} //链接标识
	data       *buffer
}

/**
 * @description: 创建一个链接
 * @param {gnet.Conn} c
 * @return {*}
 */
func NewConn(c gnet.Conn, service *Service) *Conn {
	return &Conn{ctx: context.Background(), wsConn: nil, c: c, activeTime: time.Now().UnixMicro(), service: service, data: &buffer{}}
}

/**
 * @description: 获取服务信息
 * @return {*}
 */
func (m *Conn) GetService() *Service {
	return m.service
}

/**
 * @description: 刷新活跃时间
 * @return {*}
 */
func (m *Conn) RefreshActive() {
	m.activeTime = time.Now().UnixMicro()
}

/**
 * @description: 获取活跃时间
 * @return {*}
 */
func (m *Conn) GetActiveTime() int64 {
	return m.activeTime
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
func (m *Conn) setWSConn(wsConn *websocket.Conn) {
	m.wsConn = wsConn
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
 * @description: 获取GNET链接
 * @return {*}
 */
func (m *Conn) GetConn() gnet.Conn {
	return m.c
}

/**
 * @description: 将数据写入缓存
 * @param {[]byte} buf
 * @return {*}
 */
func (m *Conn) write(buf []byte) {
	if m.data == nil {
		return
	}
	m.data.Writes(buf)
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

	ctx := context.Background()
	if m.GetService().GetOptions().IsWebsocket {
		messageType := FrameBinary
		if len(msgTypes) > 0 {
			messageType = msgTypes[0]
		}
		ctx = context.WithValue(ctx, WEBSOCKET_CONTEXT_TYPE_KEY, messageType)
	}
	_, buf, err := m.GetService().getDataHandle().Packet(ctx, buf)
	if err != nil {
		return -1, err
	}

	return m.c.Write(buf)
}

func (m *Conn) Close() error {
	if m.c == nil {
		return nil
	}
	return m.c.Close()
}

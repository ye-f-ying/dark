/*
 * @Author: yeying
 * @Date: 2025-06-01 18:13:46
 * @LastEditTime: 2025-06-01 18:13:53
 * @FilePath: \dark\gnet.go
 * @Description:
 */
package dark

import (
	"bufio"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/panjf2000/gnet/v2"
)

type DackServer struct {
	gnet.BuiltinEventEngine

	addr      string
	multicore bool
	eng       gnet.Engine
	connected int64

	isWebsocket bool
	upgrader    websocket.Upgrader
}

func NewDackServer() *DackServer {
	return &DackServer{
		isWebsocket: true,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

/**
 * @description: 设置websocket升级
 * @param {websocket.Upgrader} upgrader
 * @return {*}
 */
func (m *DackServer) SetUpgrader(upgrader websocket.Upgrader) {
	m.upgrader = upgrader
}

/**
 * @description:  实现gnet.BuiltinEventEngine OnBoot
 * @param {gnet.Engine} eng
 * @return {*}
 */
func (m *DackServer) OnBoot(eng gnet.Engine) gnet.Action {
	m.eng = eng
	return gnet.None
}

/**
 * @description:  实现gnet.BuiltinEventEngine OnOpen
 * @param {gnet.Conn} c
 * @return {*}
 */
func (m *DackServer) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetContext(NewConn(c))
	atomic.AddInt64(&m.connected, 1)
	return nil, gnet.None
}

/**
 * @description:  实现gnet.BuiltinEventEngine OnClose
 * @param {gnet.Conn} c
 * @param {error} err
 * @return {*}
 */
func (m *DackServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	atomic.AddInt64(&m.connected, -1)
	return gnet.None
}

/**
 * @description: 实现gnet.BuiltinEventEngine OnTraffic
 * @param {gnet.Conn} c
 * @return {*}
 */
func (m *DackServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	conn, ok := c.Context().(*Conn)
	if !ok {
		return gnet.Close
	}
	if m.isWebsocket { //判断是否是wensocket服务
		if conn.wsConn == nil { //是否已经升级为WEBSOCKET服务了
			_, buf, err := conn.Read()
			if err != nil {
				return gnet.Close
			}

			reader := strings.NewReader(string(buf))
			bufReader := bufio.NewReader(reader)
			req, err := http.ReadRequest(bufReader)
			if err != nil {
				return gnet.Close
			}
			w := &ResponseWriter{conn: c, header: req.Header}
			wsConn, err := m.upgrader.Upgrade(w, req, nil)
			if err != nil {
				return gnet.Close
			}
			conn.SetWSConn(wsConn)
		} else {
			messageType, buf, _ := conn.Read()
			conn.Send(buf, messageType)

		}
	}

	return gnet.None
}

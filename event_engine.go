package dark

import (
	"bufio"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/panjf2000/gnet/v2"
)

type EventEngine struct {
	gnet.BuiltinEventEngine
	service  *Service
	upgrader websocket.Upgrader
	handle   HandleInterface
}

func NewEventEngine() *EventEngine {
	return &EventEngine{service: newService(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}}
}

func (m *EventEngine) setHandle(handle HandleInterface) {
	m.handle = handle
}

/**
 * @description: 设置websocket升级
 * @param {websocket.Upgrader} upgrader
 * @return {*}
 */
func (m *EventEngine) setUpgrader(upgrader websocket.Upgrader) {
	m.upgrader = upgrader
}

/**
 * @description: 获取服务信息
 * @return {*}
 */
func (m *EventEngine) GetService() *Service {
	if m.service == nil {
		m.service = newService()
	}
	return m.service
}

/**
 * @description:  实现gnet.BuiltinEventEngine OnBoot
 * @param {gnet.Engine} eng
 * @return {*}
 */
func (m *EventEngine) OnBoot(eng gnet.Engine) gnet.Action {
	m.GetService().setEngine(eng)
	return gnet.None
}

/**
 * @description:  实现gnet.BuiltinEventEngine OnOpen
 * @param {gnet.Conn} c
 * @return {*}
 */
func (m *EventEngine) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	conn := NewConn(c)
	conn.connTime = time.Now().UnixMilli()
	c.SetContext(conn)
	m.GetService().AddConn(c, conn)
	if !m.GetService().GetOptions().IsWebsocket && m.handle != nil {
		Go(func() {
			m.handle.OnConnect(NewSession(conn, nil, 0))
		})
	}
	return nil, gnet.None
}

/**
 * @description:  实现gnet.BuiltinEventEngine OnClose
 * @param {gnet.Conn} c
 * @param {error} err
 * @return {*}
 */
func (m *EventEngine) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if m.handle != nil {
		conn, ok := c.Context().(*Conn)
		Go(func() {
			if ok {
				m.handle.OnClose(NewSession(conn, nil, 0))
			} else {
				m.handle.OnClose(NewSession(m.GetService().GetConn(c), nil, 0))
			}
		})

	}
	m.GetService().DeleteConn(c)
	return gnet.None
}

/**
 * @description: 实现gnet.BuiltinEventEngine OnTraffic
 * @param {gnet.Conn} c
 * @return {*}
 */
func (m *EventEngine) OnTraffic(c gnet.Conn) (action gnet.Action) {
	conn, ok := c.Context().(*Conn)
	if !ok {
		return gnet.Close
	}
	conn.RefreshActive()
	session := NewSession(conn, nil, 0)
	if m.GetService().GetOptions().IsWebsocket { //判断是否是wensocket服务
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
			conn.setWSConn(wsConn)
			Go(func() {
				if m.handle != nil {
					m.handle.OnConnect(session)
				}
			})
			return gnet.None
		}
		msgType, msg, err := conn.GetWSConn().ReadMessage()
		if err != nil {
			return gnet.Close
		}
		session.msg = msg
		session.msgType = msgType
	} else {
		_, buf, err := conn.Read()
		if err != nil {
			return gnet.Close
		}
		session.msg = buf
	}

	Go(func() {
		if m.handle != nil {
			m.handle.OnMessage(session)
		}
	})

	return gnet.None
}

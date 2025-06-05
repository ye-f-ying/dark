/*
 * @Author: yeying
 * @Date: 2025-06-05 21:40:00
 * @LastEditTime: 2025-06-05 22:52:08
 * @FilePath: \dark\server.go
 * @Description:
 */
package dark

import (
	"fmt"

	"github.com/panjf2000/gnet/v2"
)

type Server struct {
	engine *EventEngine
}

func NewServer() *Server {
	return &Server{engine: NewEventEngine()}
}

func (m *Server) getEngine() *EventEngine {
	if m.engine == nil {
		m.engine = NewEventEngine()
	}
	return m.engine
}

func (m *Server) Start(addr string, os ...interface{}) error {
	var opts []Option
	var g_opts []gnet.Option
	for _, v := range os {
		switch opt := (v).(type) {
		case Option:
			opts = append(opts, opt)
		case gnet.Option:
			g_opts = append(g_opts, opt)
		}
	}

	eng := m.getEngine()
	eng.GetService().SetOptions(opts...)
	return gnet.Run(eng, fmt.Sprintf("tcp://%s", addr), g_opts...)
}

func (m *Server) Close() {

}

func Run(handle HandleInterface, addr string, os ...interface{}) error {
	s := NewServer()
	defer s.Close()
	s.getEngine().SetHandle(handle)
	return s.Start(addr, os...)
}

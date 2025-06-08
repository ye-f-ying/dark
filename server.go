/*
 * @Author: yeying
 * @Date: 2025-06-05 21:40:00
 * @LastEditTime: 2025-06-08 18:31:36
 * @FilePath: \dark\server.go
 * @Description:
 */
package dark

import (
	"dark/logger"
	"dark/logger/dlog"
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

func (m *Server) Start(os ...interface{}) error {
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
	eng.GetService().setOptions(opts...)
	opt := eng.GetService().GetOptions()
	addr := opt.Addr
	if opt.isUpgrader {
		eng.setUpgrader(opt.Upgrader)
	}
	if opt.Addr == "" {
		addr = ":8989"
	}

	//日志设置
	cfg := logger.GetDefaultLoggerConfig()
	if opt.LogConsole != nil {
		cfg.Console = *opt.LogConsole
	}

	if opt.LogConsoleHuman != nil {
		cfg.ConsoleHuman = *opt.LogConsoleHuman
	}

	if opt.LogLevel != "" {
		cfg.Level = opt.LogLevel
	}

	if opt.LogSplitByLevel != nil {
		cfg.SplitByLevel = *opt.LogSplitByLevel
	}

	if opt.LogDir != "" {
		cfg.LogDir = opt.LogDir
	}

	log, err := logger.NewLogger(cfg)
	if err != nil {
		return err
	}
	dlog.SetLogger(log)
	return gnet.Run(eng, fmt.Sprintf("tcp://%s", addr), g_opts...)
}

func (m *Server) Close() {

}

func Run(handle HandleInterface, os ...interface{}) error {
	s := NewServer()
	defer s.Close()
	s.getEngine().setHandle(handle)
	return s.Start(os...)
}

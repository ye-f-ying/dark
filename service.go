/*
 * @Author: yeying
 * @Date: 2025-06-03 22:44:01
 * @LastEditTime: 2025-06-23 22:21:28
 * @FilePath: \dark\service.go
 * @Description:
 */
package dark

import (
	"github.com/panjf2000/gnet/v2"
)

type Service struct {
	opt *Options //配置信息
	eng gnet.Engine
	dh  IDataHandle
}

func newService() *Service {
	return &Service{}
}

func (m *Service) GetOptions() *Options {
	return m.opt
}

func (m *Service) setEngine(eng gnet.Engine) {
	m.eng = eng
}

func (m *Service) GetEngine(eng gnet.Engine) gnet.Engine {
	return m.eng
}

func (m *Service) getDataHandle() IDataHandle {
	if m.dh == nil {
		if m.GetOptions().IsWebsocket {
			m.dh = &WSDataHandle{}
		} else {
			m.dh = &DefaultDataHandle{}
		}
	}
	return m.dh
}

/**
 * @description:
 * @param {...Option} options
 * @return {*}
 */
func (m *Service) setOptions(options ...Option) {
	m.opt = loadOptions(options...)
}

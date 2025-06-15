/*
 * @Author: yeying
 * @Date: 2025-06-03 22:44:01
 * @LastEditTime: 2025-06-15 15:17:22
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

/**
 * @description:
 * @param {...Option} options
 * @return {*}
 */
func (m *Service) setOptions(options ...Option) {
	m.opt = loadOptions(options...)
}

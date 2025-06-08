package dark

import (
	"sync"

	"github.com/panjf2000/gnet/v2"
)

type Service struct {
	opt *Options

	eng gnet.Engine

	conns sync.Map
}

func newService() *Service {
	return &Service{}
}

/**
 * @description: 获取Options
 * @return {*}
 */
func (m *Service) GetOptions() *Options {
	return m.opt
}

/**
 * @description: 添加链接
 * @param {gnet.Conn} key
 * @param {*Conn} val
 * @return {*}
 */
func (m *Service) AddConn(key gnet.Conn, val *Conn) {
	m.conns.Store(key, val)
}

/**
 * @description: 获取链接
 * @param {gnet.Conn} key
 * @return {*}
 */
func (m *Service) GetConn(key gnet.Conn) *Conn {
	value, ok := m.conns.Load(key)
	if ok {
		if v, ok := value.(*Conn); ok {
			return v
		}
	}
	return nil
}

/**
 * @description: 删除链接
 * @param {gnet.Conn} key
 * @return {*}
 */
func (m *Service) DeleteConn(key gnet.Conn) {
	m.conns.Delete(key)
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

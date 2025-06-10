/*
 * @Author: yeying
 * @Date: 2025-06-10 21:27:18
 * @LastEditTime: 2025-06-10 22:11:38
 * @FilePath: \dark\go_pool.go
 * @Description:
 */
package dark

import (
	"github.com/panjf2000/ants/v2"
)

var goPool IGoPool

func setGoPool(p IGoPool) {
	goPool = p
}

func getGoPool() IGoPool {
	if goPool == nil {
		goPool = &DefaultPool{}
		goPool.Init()
	}
	return goPool
}

type IGoPool interface {
	Init() error
	Go(func()) error
	Close()
}
type DefaultPool struct {
	IGoPool
}

func (m *DefaultPool) Init() error {
	return nil
}

func (m *DefaultPool) Go(task func()) error {
	if task != nil {
		task()
	}
	return nil
}

func (m *DefaultPool) Close() {

}

type APool struct {
	IGoPool
	ant  *ants.Pool
	opts []ants.Option
	size int
}

func newAPool(size int, opt ...ants.Option) *APool {
	return &APool{size: size, opts: opt}
}

func (m *APool) Init() error {
	if m.size <= 0 {
		m.size = ants.DefaultAntsPoolSize
	}
	ant, err := ants.NewPool(m.size, m.opts...)
	if err != nil {
		return err
	}
	m.ant = ant
	return nil
}

func (m *APool) Go(task func()) error {
	if m.ant != nil && task != nil {
		return m.ant.Submit(task)
	}
	return &DarkError{err: ErrAntsIsNil}
}

func (m *APool) Close() {
	if m.ant != nil {
		m.ant.Release()
	}
}

/**
 * @description: 运行异步任务
 * @return {*}
 */
func Go(task func()) {
	getGoPool().Go(task)
}

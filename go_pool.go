/*
 * @Author: yeying
 * @Date: 2025-06-10 21:27:18
 * @LastEditTime: 2025-06-15 14:43:20
 * @FilePath: \dark\go_pool.go
 * @Description:
 */
package dark

import (
	"github.com/panjf2000/ants/v2"
)

// 全局协程池对象
var goPool IGoPool

/**
 * @description: 设置协程池
 * @param {IGoPool} p
 * @return {*}
 */
func setGoPool(p IGoPool) {
	goPool = p
}

/**
 * @description: 获取协程池
 * @return {*}
 */
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

// 默认协程池处理
type DefaultPool struct {
	IGoPool
}

/**
 * @description: 实现初始化
 * @return {*}
 */
func (m *DefaultPool) Init() error {
	return nil
}

/**
 * @description: 实现协程池执行
 * @return {*}
 */
func (m *DefaultPool) Go(task func()) error {
	if task != nil {
		task()
	}
	return nil
}

/**
 * @description: 实现协程池关闭
 * @return {*}
 */
func (m *DefaultPool) Close() {

}

/**
 * @description: ants协程池对象
 * @return {*}
 */
type APool struct {
	IGoPool
	ant  *ants.Pool
	opts []ants.Option
	size int
}

/**
 * @description: ants协程池
 * @param {int} size
 * @param {...ants.Option} opt
 * @return {*}
 */
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

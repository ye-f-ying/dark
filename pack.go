/*
 * @Author: yeying
 * @Date: 2025-06-23 21:41:22
 * @LastEditTime: 2025-06-23 22:09:06
 * @FilePath: \dark\pack.go
 * @Description:
 */
package dark

import (
	"context"
	"fmt"
	"sync"
)

type IPack interface {
	GetHeader(header interface{})
	GetData() []byte
	GetType() int
}
type IDataHandle interface {
	Unpacking(context.Context, []byte) (int, IPack, error)
	Packet(context.Context, []byte) (int, []byte, error)
}

type IDataHeader interface {
	GetHeader([]byte, interface{})
}

type buffer struct {
	buf  []byte // 数据缓冲区
	lock sync.Mutex
}

func (m *buffer) Writes(buf []byte) {
	if buf == nil {
		return
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	m.buf = append(m.buf, buf...)
}

func (m *buffer) Reads(ctx context.Context, h IDataHandle) ([]IPack, error) {
	if h == nil {
		return nil, fmt.Errorf("s")
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	isRead := true
	var packs []IPack
	for isRead {
		l, pack, err := h.Unpacking(ctx, m.buf)
		if err != nil {
			isRead = false
			return packs, err
		}
		if len(m.buf) < l || pack == nil {
			isRead = false
			return packs, nil
		}
		packs = append(packs, pack)
		m.buf = m.buf[l:]
		if len(m.buf) <= 0 {
			isRead = false
		}
	}
	return packs, nil
}

type DefaultDataHandle struct {
}

func (m *DefaultDataHandle) Unpacking(context.Context, []byte) (int, IPack, error) {
	return 0, nil, nil
}
func (m *DefaultDataHandle) Packet(context.Context, []byte) (int, []byte, error) {
	return 0, nil, nil
}

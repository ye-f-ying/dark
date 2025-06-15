package dark

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gammazero/deque"
	"github.com/panjf2000/gnet/v2"
)

type ConnStatus int

const (
	CONN_STATUS_NORMAL        ConnStatus = iota + 1 //正常链接 Disconnection
	CONN_STATUS_CONNECTING                          //链接中
	CONN_STATUS_DISCONNECTION                       //断线 Connecting
)

// normal
type Message struct {
	Data      []byte    //需要发送的消息
	Timestamp time.Time // 加入缓存的时间
}

// 连接状态变更事件
type ConnectionEvent struct {
	Mark     interface{}
	OldState ConnStatus
	NewState ConnStatus
	Time     time.Time
}

// 事件回调函数
type EventHandler func(event ConnectionEvent)

type connInfo struct {
	Conn     *Conn
	Mark     interface{} //链接标识
	CacheMsg deque.Deque[Message]
	Status   ConnStatus
}

type ConnManage struct {
	// TODO: 后面可以优化为数据分片，
	//链接列表
	connMap sync.Map

	//标识列表
	markMap sync.Map

	//断线标识列表
	disconnMarkMap sync.Map

	//重试次数
	withRetry int

	eventHandlers []EventHandler
}

// 注册事件回调
func (m *ConnManage) RegisterEventHandler(handler EventHandler) {
	m.eventHandlers = append(m.eventHandlers, handler)
}

// 触发事件
func (m *ConnManage) fireEvent(event ConnectionEvent) {
	for _, handler := range m.eventHandlers {
		handler(event)
	}
}

/**
 * @description: 设置重试次数
 * @param {int} retry
 * @return {*}
 */
func (m *ConnManage) SetwithRetry(retry int) {
	if retry <= 0 {
		retry = 1
	}
	m.withRetry = retry
}

/**
 * @description: 获取重试次数
 * @return {*}
 */
func (m *ConnManage) GetWithRetry() int {
	if m.withRetry <= 0 { //重试次数不能少于1次
		m.withRetry = 1
	}
	return m.withRetry
}

/**
 * @description:一个链接加入链接管理器
 * @param {gnet.Conn} conn
 * @param {interface{}} mark
 * @return {*}
 */
func (m *ConnManage) Add(conn *Conn, mark interface{}) {
	if conn == nil || mark == nil {
		return
	}
	mInfo := m.connFindByMark(mark)
	cInfo := m.connFindByConn(conn.GetConn())
	var info *connInfo
	if mInfo != nil && cInfo == nil { //表示是断线的链接
		info = mInfo
	} else if mInfo == nil && cInfo != nil { // 异常的链接
		info = cInfo
	} else {
		info = &connInfo{}
	}
	info.Conn = conn
	info.Mark = mark
	//触发事件回调
	if info.Status != CONN_STATUS_NORMAL {
		oldState := info.Status
		info.Status = CONN_STATUS_NORMAL
		m.fireEvent(ConnectionEvent{
			Mark:     mark,
			OldState: oldState,
			NewState: CONN_STATUS_NORMAL,
			Time:     time.Now(),
		})
	}
	m.connMap.Store(conn.GetConn(), info)
	m.markMap.Store(mark, info)

	if info.CacheMsg.Len() > 0 { //需要发送重连消息
		//异步发送数据
		go func() {
			for info.CacheMsg.Len() > 0 { //后面在增加间隔发送
				cacheMsg := cInfo.CacheMsg.PopFront()
				info.Conn.Send(cacheMsg.Data)
			}
		}()

	}
}

/**
 * @description: 断开链接处理
 * @param {*Conn} conn
 * @param {bool} isNormal
 * @return {*}
 */
func (m *ConnManage) HandleDisconnect(conn *Conn, isNormal bool) {
	if conn == nil {
		return
	}
	cInfo := m.connFindByConn(conn.GetConn())
	if cInfo == nil {
		return
	}
	m.connMap.Delete(conn.GetConn())
	if isNormal { //正常断开
		m.markMap.Delete(cInfo.Mark)
		return
	}
	cInfo.Conn = nil //异常断开的conn 不能在使用了
	//cInfo.Status = CONN_STATUS_DISCONNECTION
	//触发事件回调
	if cInfo.Status != CONN_STATUS_DISCONNECTION {
		oldState := cInfo.Status
		cInfo.Status = CONN_STATUS_DISCONNECTION
		m.fireEvent(ConnectionEvent{
			Mark:     cInfo.Mark,
			OldState: oldState,
			NewState: CONN_STATUS_DISCONNECTION,
			Time:     time.Now(),
		})
	}
	//加入断线列表
	m.disconnMarkMap.Store(cInfo.Mark, cInfo)
}

/**
 * @description:发送消息
 * @param {interface{}} mark
 * @param {[]byte} msg
 * @return {*}
 */
func (m *ConnManage) SendMessage(mark interface{}, msg []byte) (int, error) {
	if len(msg) <= 0 || mark == nil {
		return 0, fmt.Errorf("invalid parameters")
	}

	cInfo := m.connFindByMark(mark)
	if cInfo == nil {
		return 0, fmt.Errorf("connection not found for mark: %v", mark)
	}

	// 使用带上下文的发送，支持超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if cInfo.Conn == nil || cInfo.Status == CONN_STATUS_DISCONNECTION {
		m.cacheMessage(cInfo, msg)
		return len(msg), fmt.Errorf("connection is disconnected, message cached")
	}

	// 使用带重试的发送
	return m.sendWithRetry(ctx, cInfo, msg)
}

func (m *ConnManage) sendWithRetry(ctx context.Context, cInfo *connInfo, msg []byte) (int, error) {
	var lastErr error
	withRetry := m.GetWithRetry()
	for i := 0; i < withRetry; i++ {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			n, err := cInfo.Conn.Send(msg)
			if err == nil {
				return n, nil
			}
			lastErr = err
			time.Sleep(time.Millisecond * 100 * time.Duration(i+1))
		}
	}
	return 0, fmt.Errorf("failed after retries: %w", lastErr)
}

/**
 * @description:
 * @param {*connInfo} cInfo
 * @param {[]byte} data
 * @return {*}
 */
func (m *ConnManage) cacheMessage(cInfo *connInfo, data []byte) {
	if cInfo == nil || len(data) <= 0 {
		return
	}
	if cInfo.CacheMsg.Len() > 100 {
		cInfo.CacheMsg.PopFront()
	}
	cInfo.CacheMsg.PushBack(Message{
		Data: data, Timestamp: time.Now(),
	})
}

func (m *ConnManage) connFindByMark(mark interface{}) *connInfo {
	if info, ok := m.markMap.Load(mark); ok {
		return info.(*connInfo)
	}
	return nil
}

func (m *ConnManage) connFindByConn(c gnet.Conn) *connInfo {
	if info, ok := m.connMap.Load(c); ok {
		return info.(*connInfo)
	}
	return nil
}

/**
 * @description: 通过标识获取链接信息
 * @param {interface{}} mark
 * @return {*}
 */
func (m *ConnManage) FindByMark(mark interface{}) *Conn {
	info := m.connFindByMark(mark)
	if info != nil {
		return info.Conn
	}
	return nil
}

/**
 * @description: 通链接获取链接
 * @param {*Conn} conn
 * @return {*}
 */
func (m *ConnManage) FindByConn(conn *Conn) *Conn {
	if conn == nil {
		return nil
	}
	return m.FindByGConn(conn.GetConn())
}

/**
 * @description: 通过gnet的指针获取链接信息
 * @param {gnet.Conn} conn
 * @return {*}
 */
func (m *ConnManage) FindByGConn(conn gnet.Conn) *Conn {
	if conn == nil {
		return nil
	}
	info := m.connFindByConn(conn)
	if info != nil {
		return info.Conn
	}
	return nil
}

/**
 * @description: 启动超时检查
 * @param {time.Duration} interval
 * @return {*}
 */
func (m *ConnManage) StartHeartbeatChecker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			m.checkHeartbeats()
		}
	}()
}

/**
 * @description: 检查是否超时链接
 * @return {*}
 */
func (m *ConnManage) checkHeartbeats() {
	now := time.Now().UnixNano()
	timeout := int64(30 * time.Second) // 30秒超时

	m.markMap.Range(func(key, value interface{}) bool {
		info := value.(*connInfo)
		if info.Conn != nil {
			lastActive := info.Conn.GetActiveTime()
			if now-lastActive > timeout {
				// 连接超时，标记为断开
				m.HandleDisconnect(info.Conn, false)
			}
		}
		return true
	})
}

/**
 * @description:
 * @param {context.Context} ctx
 * @return {*}
 */
func (m *ConnManage) Shutdown(ctx context.Context) error {
	var wg sync.WaitGroup
	m.markMap.Range(func(key, value interface{}) bool {
		info := value.(*connInfo)
		if info.Conn != nil {
			wg.Add(1)
			go func(conn *Conn) {
				defer wg.Done()
				// 优雅关闭连接
				conn.Close()
			}(info.Conn)
		}
		return true
	})

	// 等待所有连接关闭
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()

	// 等待关闭完成或超时
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

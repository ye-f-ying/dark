/*
 * @Author: yeying
 * @Date: 2025-06-01 18:09:25
 * @LastEditTime: 2025-06-01 18:09:31
 * @FilePath: \dark\response_writer.go
 * @Description:
 */
package dark

import (
	"bufio"
	"net"
	"net/http"

	"github.com/panjf2000/gnet/v2"
)

type ResponseWriter struct {
	http.ResponseWriter
	http.Hijacker
	conn   gnet.Conn
	status int
	header http.Header
}

func (m *ResponseWriter) Header() http.Header {
	if m.header == nil {
		m.header = make(http.Header)
	}
	return m.header
}

func (m *ResponseWriter) WriteHeader(statusCode int) {
	m.status = statusCode
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	return w.conn.Write(data)
}

// 实现 http.Hijacker 接口
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.conn,
		bufio.NewReadWriter(bufio.NewReader(w.conn),
			bufio.NewWriter(w.conn)), nil
}

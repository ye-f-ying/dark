/*
 * @Author: yeying
 * @Date: 2025-05-28 21:49:21
 * @LastEditTime: 2025-06-08 18:33:59
 * @FilePath: \dark\examples\websocket\server\main.go
 * @Description:
 */
package main

import (
	"dark"
	"dark/logger/dlog"
	"fmt"
)

type Test struct {
	dark.DarkHandle
}

func (m *Test) OnMessage(c *dark.Session) {
	fmt.Println(c.GetMsg())
	dlog.Info(string(c.GetMsg()))
}

func main() {
	dark.Run(&Test{}, dark.WithAddr(":12000"), dark.WithIsWebsocket(true))
}

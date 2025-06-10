/*
 * @Author: yeying
 * @Date: 2025-05-28 21:49:21
 * @LastEditTime: 2025-06-10 22:02:25
 * @FilePath: \dark\examples\websocket\server\main.go
 * @Description:
 */
package main

import (
	"dark"
	"dark/logger/dlog"
	"fmt"

	"github.com/panjf2000/ants/v2"
)

type Test struct {
	dark.DarkHandle
}

func (m *Test) OnMessage(c *dark.Session) {
	fmt.Println(c.GetMsg())
	dlog.Info(string(c.GetMsg()))
}

func main() {
	dark.Run(&Test{}, dark.WithAddr(":12000"), dark.WithIsWebsocket(true), dark.WithAnts(100, ants.WithPreAlloc(true)))
}

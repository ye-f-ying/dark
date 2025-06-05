/*
 * @Author: yeying
 * @Date: 2025-05-28 21:49:21
 * @LastEditTime: 2025-06-05 22:49:35
 * @FilePath: \dark\examples\websocket\server\main.go
 * @Description:
 */
package main

import (
	"dark"
	"fmt"
)

type Test struct {
	dark.DarkHandle
}

func (m *Test) OnMessage(c *dark.Session) {
	fmt.Println(c.GetMsg())
}

func main() {
	dark.Run(&Test{}, ":12000", dark.WithIsWebsocket(true))
}

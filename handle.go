package dark

type HandleInterface interface {
	OnConnect(s *Session)
	OnMessage(s *Session)
	OnClose(s *Session)
}

type DarkHandle struct {
	HandleInterface
}

func (*DarkHandle) OnConnect(s *Session) {

}

func (*DarkHandle) OnMessage(s *Session) {

}

func (*DarkHandle) OnClose(s *Session) {

}

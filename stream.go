package main

import "github.com/uhoh-itsmaciek/femebe/core"

type loggedStream struct {
	debug  LogFunc
	stream *core.MessageStream
}

func (s loggedStream) log(dir string, m *core.Message) {
	if m.MsgType() != 0 {
		s.debug("dir", dir, "type", string(m.MsgType()), "size", m.Size())
	} else {
		b, _ := m.Force()
		if len(b) == 4 && b[0] == 0x04 && b[1] == 0xd2 && b[2] == 0x16 && b[3] == 0x2f {
			s.debug("dir", dir, "type", "SSL")
		} else {
			s.debug("dir", dir, "type", "Start")
		}
	}
}

func (s loggedStream) Next(dir string, m *core.Message) error {
	err := s.stream.Next(m)
	if err == nil {
		s.log(dir, m)
	}
	return err
}

func (s loggedStream) Send(dir string, m *core.Message) error {
	s.log(dir, m)
	return s.stream.Send(m)
}

func (s loggedStream) SendSSLRequestResponse(dir string, r byte) error {
	s.debug("dir", dir, "type", "SSL", "response", string(r))
	return s.stream.SendSSLRequestResponse(r)
}

type BackendStream loggedStream

func (be BackendStream) Close() error  { return be.stream.Close() }
func (be BackendStream) Flush() error  { return be.stream.Flush() }
func (be BackendStream) HasNext() bool { return be.stream.HasNext() }

func (be BackendStream) Next(m *core.Message) error { return (loggedStream)(be).Next(" <B", m) }
func (be BackendStream) Send(m *core.Message) error { return (loggedStream)(be).Send(" >B", m) }

type FrontendStream loggedStream

func (fe FrontendStream) Close() error  { return fe.stream.Close() }
func (fe FrontendStream) Flush() error  { return fe.stream.Flush() }
func (fe FrontendStream) HasNext() bool { return fe.stream.HasNext() }

func (fe FrontendStream) Next(m *core.Message) error { return (loggedStream)(fe).Next("F> ", m) }
func (fe FrontendStream) Send(m *core.Message) error { return (loggedStream)(fe).Send("F< ", m) }

func (fe FrontendStream) SendSSLRequestResponse(r byte) error {
	return (loggedStream)(fe).SendSSLRequestResponse("F< ", r)
}

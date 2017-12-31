package pgtwixt

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/uhoh-itsmaciek/femebe/core"
	"github.com/uhoh-itsmaciek/femebe/proto"
)

type CancellationKey struct {
	id, secret uint32
}

type Server struct {
	Debug LogFunc
	Info  LogFunc
	tls   *tls.Config

	Cancel  func(CancellationKey)
	Session func(FrontendStream, map[string]string)
}

func (s *Server) accept(conn net.Conn) {
	if err := s.handshake(conn); err != nil {
		s.Info("msg", "Error during handshake", "error", err)
	}
}

// handshake interprets the initial SSL, Startup, and/or Cancel message(s).
func (s *Server) handshake(conn net.Conn) (err error) {
	var msg core.Message
	fe := FrontendStream{
		debug:  s.Debug,
		stream: core.NewFrontendStream(conn),
	}

	if err = fe.Next(&msg); err != nil {
		return
	}

	if proto.IsSSLRequest(&msg) {
		if err = msg.Discard(); err != nil {
			return
		}
		if s.tls == nil {
			if err = fe.SendSSLRequestResponse(core.RejectSSLRequest); err != nil {
				return
			}
		} else {
			if err = fe.SendSSLRequestResponse(core.AcceptSSLRequest); err != nil {
				return
			}

			tlsConn := tls.Server(conn, s.tls)
			if err = tlsConn.Handshake(); err != nil {
				return
			}

			fe.stream = core.NewFrontendStream(tlsConn)
		}
		if err = fe.Next(&msg); err != nil {
			return
		}
	}

	if proto.IsStartupMessage(&msg) {
		var su *proto.StartupMessage
		if su, err = proto.ReadStartupMessage(&msg); err == nil {
			s.Session(fe, su.Params)
		}
		return
	}

	if proto.IsCancelRequest(&msg) {
		var c *proto.CancelRequest
		if c, err = proto.ReadCancelRequest(&msg); err == nil {
			s.Cancel(CancellationKey{
				id:     c.BackendPid,
				secret: c.SecretKey,
			})
		}
		_ = fe.Close()
		return
	}

	return fmt.Errorf("Unknown message: %q", msg.MsgType())
}

func (s *Server) Serve(l net.Listener) error {
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				// TODO back off
				time.Sleep(5 * time.Millisecond)
				continue
			}
			return err
		}

		go s.accept(conn)
	}
}

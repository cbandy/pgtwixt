package pgtwixt

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"time"

	"github.com/uhoh-itsmaciek/femebe/core"
	"github.com/uhoh-itsmaciek/femebe/proto"
	"github.com/uhoh-itsmaciek/femebe/util"
)

type Connector struct {
	Dialer
}

// Cancel opens a new connection to the backend and sends a CancelRequest.
func (cn Connector) Cancel(c CancellationKey) error {
	be, err := cn.Dial(context.Background())

	if err == nil {
		var msg core.Message
		proto.InitCancelRequest(&msg, c.id, c.secret)

		err = be.Send(&msg)
		_ = be.Close()
	}

	return err
}

// Startup opens a new connection to the backend and sends a StartupMessage.
func (cn Connector) Startup(options map[string]string) (BackendStream, error) {
	be, err := cn.Dial(context.Background())

	if err == nil {
		var msg core.Message
		proto.InitStartupMessage(&msg, options)

		err = be.Send(&msg)
		if err != nil {
			_ = be.Close()
		}
	}

	return be, err
}

type Dialer interface {
	Dial(context.Context) (BackendStream, error)
}

type TCPDialer struct {
	Debug LogFunc

	Addr      string // "yahoo.com:8080" "1.2.3.4:9999"
	SSLMode   string
	SSLConfig tls.Config
	Timeout   time.Duration

	KeepAlivesCount    int
	KeepAlivesDisable  bool
	KeepAlivesIdle     time.Duration
	KeepAlivesInterval time.Duration
	// https://github.com/golang/go/blob/master/src/net/tcpsockopt_*.go
}

// Dial opens a new connection to the backend, negotiates any TLS upgrade, and verifies server certificates.
func (d TCPDialer) Dial(ctx context.Context) (BackendStream, error) {
	nd := net.Dialer{Timeout: d.Timeout}
	conn, err := nd.DialContext(ctx, "tcp", d.Addr)

	if err == nil {
		err = conn.(*net.TCPConn).SetKeepAlive(true)
	}
	if err == nil {
		cfg := util.SSLConfig{Mode: util.SSLMode(d.SSLMode), Config: d.SSLConfig}
		cfg.Config.InsecureSkipVerify = d.SSLMode != "verify-full"
		conn, err = util.NegotiateTLS(conn, &cfg)
	}
	if err == nil {
		err = d.verify(conn)
	}

	return BackendStream{
		debug:  d.Debug,
		stream: core.NewBackendStream(conn),
	}, err
}

func (d TCPDialer) verify(conn net.Conn) error {
	tc, ok := conn.(*tls.Conn)
	if !ok {
		return nil
	}

	err := tc.Handshake()
	if err == nil && d.SSLMode == "verify-ca" {
		s := tc.ConnectionState()
		v := x509.VerifyOptions{Roots: d.SSLConfig.RootCAs}

		if len(s.PeerCertificates) > 0 {
			_, err = s.PeerCertificates[0].Verify(v)
		}
	}

	return err
}

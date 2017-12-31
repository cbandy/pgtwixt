package pgtwixt

import (
	"github.com/uhoh-itsmaciek/femebe/core"
	"github.com/uhoh-itsmaciek/femebe/proto"
	"github.com/uhoh-itsmaciek/femebe/util"
)

type Connector struct {
	Debug    LogFunc
	Dialer   Dialer
	Location string
}

// Cancel opens a new connection to the backend and sends a CancelRequest.
func (cn Connector) Cancel(c CancellationKey) error {
	be, err := cn.Dialer.Dial(cn.Debug, cn.Location)

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
	be, err := cn.Dialer.Dial(cn.Debug, cn.Location)

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

type Dialer util.SSLConfig

// Dial opens a new connection to the backend and negotiates any TLS upgrade.
func (d Dialer) Dial(debug LogFunc, location string) (BackendStream, error) {
	conn, err := util.AutoDial(location)

	if err == nil {
		conn, err = util.NegotiateTLS(conn, (*util.SSLConfig)(&d))
	}

	return BackendStream{
		debug:  debug,
		stream: core.NewBackendStream(conn),
	}, err
}

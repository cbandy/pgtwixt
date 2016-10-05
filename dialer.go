package main

import (
	"github.com/uhoh-itsmaciek/femebe/core"
	"github.com/uhoh-itsmaciek/femebe/util"
)

type Dialer util.SSLConfig

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

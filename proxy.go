package pgtwixt

import (
	"io"

	"github.com/uhoh-itsmaciek/femebe/core"
)

type Proxy struct {
	Info    LogFunc
	Startup func(map[string]string) (BackendStream, error)
}

// pump copies messages one way between two streams, blocking and flushing when
// necessary.
func (*Proxy) pump(errc chan<- error, from, to core.Stream) {
	var err error
	var msg core.Message

	for {
		if err = from.Next(&msg); err != nil {
			break
		}
		if err = to.Send(&msg); err != nil {
			break
		}
		if !from.HasNext() {
			if err = to.Flush(); err != nil {
				break
			}
		}
	}

	errc <- err
}

func (p *Proxy) Run(fe FrontendStream, startup map[string]string) {
	defer fe.Close()

	be, err := p.Startup(startup)
	if err != nil {
		p.Info("msg", "Error connecting to backend", "error", err)
		return
	}
	defer be.Close()

	errc := make(chan error, 1)
	go p.pump(errc, be, fe)
	go p.pump(errc, fe, be)

	err = <-errc
	if err != nil && err != io.EOF {
		p.Info("msg", "Error while proxying", "error", err)
		return
	}
}

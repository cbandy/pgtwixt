package pgtwixt

import (
	"bytes"
	"testing"

	"github.com/uhoh-itsmaciek/femebe/core"
	"github.com/uhoh-itsmaciek/femebe/proto"

	"github.com/stretchr/testify/assert"
)

func TestServerAcceptCancel(t *testing.T) {
	t.Parallel()

	var msg core.Message
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	conn := bufConn{nopCloser{buf}}
	srv := Server{
		Debug: func(...interface{}) error { return nil },
	}

	testCancel := func(t *testing.T) {
		proto.InitCancelRequest(&msg, 2600, 1957)
		msg.WriteTo(buf)

		var called bool
		srv.Cancel = func(c CancellationKey) {
			called = true

			assert.Equal(t, CancellationKey{id: 2600, secret: 1957}, c)
		}
		srv.accept(conn)

		assert.True(t, called, "Expected cancel to be called")
	}

	testStartup := func(t *testing.T) {
		proto.InitStartupMessage(&msg, map[string]string{"user": "mary"})
		msg.WriteTo(buf)

		var called bool
		srv.Session = func(fe FrontendStream, startup map[string]string) {
			called = true

			assert.Equal(t, map[string]string{"user": "mary"}, startup)
			assert.NotNil(t, fe.debug, "Expected stream to have logger, got none")
		}
		srv.accept(conn)

		assert.True(t, called, "Expected session to be called")
	}

	buf.Reset()
	t.Run("Cancel", testCancel)

	buf.Reset()
	t.Run("Startup", testStartup)

	buf.Reset()
	buf.Write(sslRequest)
	t.Run("SSL,Cancel", testCancel)

	buf.Reset()
	buf.Write(sslRequest)
	t.Run("SSL,Startup", testStartup)
}

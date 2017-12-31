package pgtwixt

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/uhoh-itsmaciek/femebe/core"
	"github.com/uhoh-itsmaciek/femebe/proto"
)

func TestStreamInterfaces(t *testing.T) {
	var _ core.Stream = BackendStream{}
	var _ core.Stream = FrontendStream{}
}

func TestBackendStreamLogs(t *testing.T) {
	t.Parallel()

	var logged []interface{}
	var msg core.Message

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	be := BackendStream{
		stream: core.NewBackendStream(nopCloser{buf}),
		debug:  func(keyvals ...interface{}) error { logged = keyvals; return nil },
	}

	t.Run("Send", func(t *testing.T) {
		msg.InitFromBytes('B', []byte{0, 0, 0, 0})
		err := be.Send(&msg)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if expected := []interface{}{"dir", " >B", "type", "B", "size", uint32(8)}; !reflect.DeepEqual(logged, expected) {
			t.Errorf("Expected %#v, got %#v", expected, logged)
		}
	})

	buf.Reset()

	t.Run("Next", func(t *testing.T) {
		msg.InitFromBytes('B', []byte{0, 0, 0, 0})
		msg.WriteTo(buf)
		err := be.Next(&msg)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if expected := []interface{}{"dir", " <B", "type", "B", "size", uint32(8)}; !reflect.DeepEqual(logged, expected) {
			t.Errorf("Expected %#v, got %#v", expected, logged)
		}
	})
}

func BenchmarkBackendStreamSend(b *testing.B) {
	var msg core.Message
	be := BackendStream{
		stream: core.NewBackendStream(nopReaderCloser{ioutil.Discard}),
		debug:  func(...interface{}) error { return nil },
	}

	for i := 0; i < b.N; i++ {
		be.Send(&msg)
	}
}

func TestFrontendStreamLogs(t *testing.T) {
	t.Parallel()

	var logged []interface{}
	var msg core.Message

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	fe := FrontendStream{
		stream: core.NewFrontendStream(nopCloser{buf}),
		debug:  func(keyvals ...interface{}) error { logged = keyvals; return nil },
	}

	t.Run("Send,SSL", func(t *testing.T) {
		err := fe.SendSSLRequestResponse('S')

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if expected := []interface{}{"dir", "F< ", "type", "SSL", "response", "S"}; !reflect.DeepEqual(logged, expected) {
			t.Errorf("Expected %#v, got %#v", expected, logged)
		}
	})

	buf.Reset()

	t.Run("Send,Start", func(t *testing.T) {
		err := fe.Send(&msg)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if expected := []interface{}{"dir", "F< ", "type", "Start"}; !reflect.DeepEqual(logged, expected) {
			t.Errorf("Expected %#v, got %#v", expected, logged)
		}
	})

	buf.Reset()

	t.Run("Next,SSL", func(t *testing.T) {
		buf.Write(sslRequest)
		err := fe.Next(&msg)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if expected := []interface{}{"dir", "F> ", "type", "SSL"}; !reflect.DeepEqual(logged, expected) {
			t.Errorf("Expected %#v, got %#v", expected, logged)
		}
	})

	buf.Reset()

	t.Run("Next,Start", func(t *testing.T) {
		proto.InitStartupMessage(&msg, nil)
		msg.WriteTo(buf)
		err := fe.Next(&msg)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if expected := []interface{}{"dir", "F> ", "type", "Start"}; !reflect.DeepEqual(logged, expected) {
			t.Errorf("Expected %#v, got %#v", expected, logged)
		}
	})
}

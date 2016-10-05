package main

import (
	"io"
	"net"
	"time"
)

var sslRequest = []byte{
	0x00, 0x00, 0x00, 0x08,
	0x04, 0xd2, 0x16, 0x2f,
}

type bufConn struct{ io.ReadWriteCloser }

func (bufConn) LocalAddr() net.Addr  { return &net.UnixAddr{} }
func (bufConn) RemoteAddr() net.Addr { return &net.UnixAddr{} }

func (bufConn) SetDeadline(time.Time) error      { return nil }
func (bufConn) SetReadDeadline(time.Time) error  { return nil }
func (bufConn) SetWriteDeadline(time.Time) error { return nil }

type nopCloser struct{ io.ReadWriter }

func (nopCloser) Close() error { return nil }

type nopReaderCloser struct{ io.Writer }

func (nopReaderCloser) Read([]byte) (int, error) { return 0, nil }
func (nopReaderCloser) Close() error             { return nil }

package pgtwixt

import (
	"context"
	"fmt"
	"net"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/uhoh-itsmaciek/femebe/core"
)

type UnixDialer struct {
	Debug LogFunc

	Addr        string // "/var/run/postgresql/.s.PGSQL.5432"
	RequirePeer string
	Timeout     time.Duration
}

// Dial opens a new connection to the backend and verifies the owner of the socket.
func (d *UnixDialer) Dial(ctx context.Context) (BackendStream, error) {
	nd := net.Dialer{Timeout: d.Timeout}
	conn, err := nd.DialContext(ctx, "unix", d.Addr)

	if err == nil {
		err = d.verify(conn)
	}

	return BackendStream{
		debug:  d.Debug,
		stream: core.NewBackendStream(conn),
	}, err
}

func (d *UnixDialer) verify(conn net.Conn) error {
	if d.RequirePeer == "" {
		return nil
	}

	var (
		ucred *syscall.Ucred
		ucerr error
	)

	// https://github.com/golang/go/issues/22953
	src, err := conn.(*net.UnixConn).SyscallConn()
	if err == nil {
		err = src.Control(func(fd uintptr) {
			ucred, ucerr = syscall.GetsockoptUcred(int(fd),
				syscall.SOL_SOCKET, syscall.SO_PEERCRED)
		})
	}
	if err == nil {
		err = ucerr
	}

	if err == nil {
		var u *user.User
		u, err = user.LookupId(strconv.Itoa(int(ucred.Uid)))
		if err == nil && u.Username != d.RequirePeer {
			err = fmt.Errorf("peer user name %q is not %q", u.Username, d.RequirePeer)
		}
	}

	return err
}

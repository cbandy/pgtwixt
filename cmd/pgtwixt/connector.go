package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/cbandy/pgtwixt"
)

func dialers(cs pgtwixt.ConnectionString) ([]pgtwixt.Dialer, error) {
	if len(cs.Host) > 0 && len(cs.HostAddr) > 0 && len(cs.Host) != len(cs.HostAddr) {
		return nil, fmt.Errorf("host and address lengths do not match: %v versus %v", cs.Host, cs.HostAddr)
	}
	if len(cs.Host) > 0 && len(cs.Port) > 1 && len(cs.Host) != len(cs.Port) {
		return nil, fmt.Errorf("host and port lengths do not match: %v versus %v", cs.Host, cs.Port)
	}
	if len(cs.HostAddr) > 0 && len(cs.Port) > 1 && len(cs.HostAddr) != len(cs.Port) {
		return nil, fmt.Errorf("address and port lengths do not match: %v versus %v", cs.HostAddr, cs.Port)
	}

	get := func(ss []string, i int, d string) string {
		if len(ss) > i {
			return ss[i]
		}
		return d
	}

	var ds []pgtwixt.Dialer
	var err error

	for i, more := 0, true; more && err == nil; more = (i < len(cs.Host) || i < len(cs.HostAddr)) {
		ds = append(ds, nil)
		port := get(cs.Port, i, get(cs.Port, 0, "5432"))

		if len(cs.HostAddr) > i || (len(cs.Host) > i && !strings.HasPrefix(cs.Host[i], "/")) {
			host := get(cs.Host, i, "")
			addr := get(cs.HostAddr, i, "")

			ds[i], err = tcpDialer(host, addr, port, cs)
		} else {
			host := get(cs.Host, i, "/tmp")

			ds[i], err = unixDialer(host, port, cs)
		}

		i++
	}

	return ds, err
}

func tcpDialer(host, hostaddr, port string, cs pgtwixt.ConnectionString) (pgtwixt.TCPDialer, error) {
	var d pgtwixt.TCPDialer
	var err error

	if hostaddr != "" {
		d.Addr = net.JoinHostPort(hostaddr, port)
	} else {
		d.Addr = net.JoinHostPort(host, port)
	}

	d.SSLMode = cs.SSLMode
	d.SSLConfig = tls.Config{
		MinVersion:    tls.VersionTLS12,
		Renegotiation: tls.RenegotiateFreelyAsClient,
		ServerName:    host,
	}

	if cs.ConnectTimeout != "" {
		d.Timeout, err = cs.SecondsDuration(cs.ConnectTimeout)
	}

	return d, err
}

func unixDialer(host, port string, cs pgtwixt.ConnectionString) (pgtwixt.UnixDialer, error) {
	var d pgtwixt.UnixDialer
	var err error

	d.Addr = host + "/.s.PGSQL." + port
	d.RequirePeer = cs.RequirePeer

	if cs.ConnectTimeout != "" {
		d.Timeout, err = cs.SecondsDuration(cs.ConnectTimeout)
	}

	return d, err
}

package pgtwixt

import (
	"strconv"
	"strings"
	"time"
)

type ConnectionString struct {
	Host, HostAddr, Port []string

	Database        string // dbname
	User, Password  string
	PasswordPath    string // passfile
	ConnectTimeout  string // connect_timeout
	ClientEncoding  string // client_encoding
	Options         string
	ApplicationName string // application_name
	FallbackName    string // fallback_application_name

	KeepAlives         string // keepalives
	KeepAlivesCount    string // keepalives_count
	KeepAlivesIdle     string // keepalives_idle
	KeepAlivesInterval string // keepalives_interval

	SSLMode        string // sslmode
	SSLRequire     string // requiressl
	SSLCompression string // sslcompression
	SSLCertPath    string // sslcert
	SSLKeyPath     string // sslkey
	SSLCAPath      string // sslrootcert
	SSLCRLPath     string // sslcrl

	RequirePeer string // requirepeer
	Service     string

	Remainder map[string]string
}

func (c *ConnectionString) Parse(s string) error {
	var g Grammar
	g.Value = func(key, value string) error {
		switch key {
		case "host":
			c.Host = strings.Split(value, ",")
		case "hostaddr":
			c.HostAddr = strings.Split(value, ",")
		case "port":
			c.Port = strings.Split(value, ",")
		case "dbname":
			c.Database = value
		case "user":
			c.User = value
		case "password":
			c.Password = value
		case "passfile":
			c.PasswordPath = value
		case "connect_timeout":
			c.ConnectTimeout = value
		case "client_encoding":
			c.ClientEncoding = value
		case "options":
			c.Options = value
		case "application_name":
			c.ApplicationName = value
		case "fallback_application_name":
			c.FallbackName = value
		case "keepalives":
			c.KeepAlives = value
		case "keepalives_idle":
			c.KeepAlivesIdle = value
		case "keepalives_interval":
			c.KeepAlivesInterval = value
		case "keepalives_count":
			c.KeepAlivesCount = value
		case "sslmode":
			c.SSLMode = value
		case "requiressl":
			c.SSLRequire = value
		case "sslcompression":
			c.SSLCompression = value
		case "sslcert":
			c.SSLCertPath = value
		case "sslkey":
			c.SSLKeyPath = value
		case "sslrootcert":
			c.SSLCAPath = value
		case "sslcrl":
			c.SSLCRLPath = value
		case "requirepeer":
			c.RequirePeer = value
		case "service":
			c.Service = value
		default:
			c.Remainder[key] = value
		}
		return nil
	}

	c.Remainder = make(map[string]string)
	return g.Parse(s)
}

func (ConnectionString) SecondsDuration(s string) (time.Duration, error) {
	n, err := strconv.Atoi(s)
	if err == nil {
		return time.Duration(n) * time.Second, nil
	}
	return 0, err
}

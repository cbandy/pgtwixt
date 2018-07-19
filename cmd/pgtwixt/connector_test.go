package main

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/cbandy/pgtwixt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDialersDefault(t *testing.T) {
	t.Parallel()

	ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{})
	require.NoError(t, err)
	require.Len(t, ds, 1)

	assert.Equal(t, pgtwixt.UnixDialer{
		Address: "/tmp/.s.PGSQL.5432",
	}, ds[0])
}

func TestDialersTCP(t *testing.T) {
	t.Parallel()

	t.Run("Host", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			Host: []string{"example.com"},
		})
		require.NoError(t, err)
		require.Len(t, ds, 1)

		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "example.com:5432",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
				ServerName:    "example.com",
			},
		}, ds[0])
	})

	t.Run("HostAddr", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			HostAddr: []string{"127.0.0.1"},
		})
		require.NoError(t, err)
		require.Len(t, ds, 1)

		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "127.0.0.1:5432",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
			},
		}, ds[0])
	})

	t.Run("HostPort", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			Host: []string{"example.com"},
			Port: []string{"888"},
		})
		require.NoError(t, err)
		require.Len(t, ds, 1)

		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "example.com:888",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
				ServerName:    "example.com",
			},
		}, ds[0])
	})

	t.Run("HostAddrPort", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			HostAddr: []string{"127.0.0.1"},
			Port:     []string{"99"},
		})
		require.NoError(t, err)
		require.Len(t, ds, 1)

		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "127.0.0.1:99",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
			},
		}, ds[0])
	})

	// TODO keepalives

	t.Run("SSL", func(t *testing.T) {
		t.Run("Mode", func(t *testing.T) {
			ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
				Host:    []string{"example.com"},
				SSLMode: "something",
			})
			require.NoError(t, err)
			require.Len(t, ds, 1)

			assert.Equal(t, pgtwixt.TCPDialer{
				Address: "example.com:5432",
				SSLMode: "something",
				SSLConfig: tls.Config{
					MinVersion:    tls.VersionTLS12,
					Renegotiation: tls.RenegotiateFreelyAsClient,
					ServerName:    "example.com",
				},
			}, ds[0])
		})

		t.Run("Require", func(t *testing.T) {
			// TODO reject
		})

		t.Run("Compression", func(t *testing.T) {
			// TODO reject
		})

		t.Run("Cert+Key", func(t *testing.T) {
			// TODO
			// - create temporary files containing certificate and key
			// - set the paths in connection string
			// - # tls.LoadX509KeyPair()
			// - assert d.SSLConfig.Certificates
		})

		t.Run("CA", func(t *testing.T) {
			// TODO
			// - create a temporary file containing certificate
			// - set the path in connection string
			// - # x509.NewCertPool()
			// - assert d.SSLConfig.RootCAs.Subjects()
		})

		t.Run("CRL", func(t *testing.T) {
			// TODO reject
		})
	})

	t.Run("Timeout", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			Host:           []string{"example.com"},
			ConnectTimeout: "10",
		})
		require.NoError(t, err)
		require.Len(t, ds, 1)

		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "example.com:5432",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
				ServerName:    "example.com",
			},
			Timeout: 10 * time.Second,
		}, ds[0])
	})
}

func TestDialersTCPError(t *testing.T) {
	t.Parallel()

	t.Run("Timeout", func(t *testing.T) {
		_, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			Host:           []string{"example.com"},
			ConnectTimeout: "nope",
		})
		assert.Error(t, err)
	})
}

func TestDialersUnix(t *testing.T) {
	t.Parallel()

	t.Run("Host", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			Host: []string{"/var/run/postgresql"},
		})
		require.NoError(t, err)
		require.Len(t, ds, 1)

		assert.Equal(t, pgtwixt.UnixDialer{
			Address: "/var/run/postgresql/.s.PGSQL.5432",
		}, ds[0])
	})

	t.Run("Port", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			Port: []string{"999"},
		})
		require.NoError(t, err)
		require.Len(t, ds, 1)

		assert.Equal(t, pgtwixt.UnixDialer{
			Address: "/tmp/.s.PGSQL.999",
		}, ds[0])
	})

	t.Run("Others", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			ConnectTimeout: "10",
			RequirePeer:    "baz",
		})
		require.NoError(t, err)
		require.Len(t, ds, 1)

		assert.Equal(t, pgtwixt.UnixDialer{
			Address:     "/tmp/.s.PGSQL.5432",
			RequirePeer: "baz",
			Timeout:     10 * time.Second,
		}, ds[0])
	})
}

func TestDialersUnixError(t *testing.T) {
	t.Parallel()

	t.Run("Timeout", func(t *testing.T) {
		_, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			ConnectTimeout: "nope",
		})
		assert.Error(t, err)
	})
}

func TestDialersMultipleError(t *testing.T) {
	t.Parallel()

	t.Run("TooManyHosts", func(t *testing.T) {
		_, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			Host:     []string{"1", "2"},
			HostAddr: []string{"1"},
		})
		assert.Error(t, err)
	})

	t.Run("TooManyHostAddr", func(t *testing.T) {
		_, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			Host:     []string{"1"},
			HostAddr: []string{"1", "2"},
		})
		assert.Error(t, err)
	})

	t.Run("TooManyPort", func(t *testing.T) {
		t.Run("Host", func(t *testing.T) {
			_, err := Connector{}.Dialers(pgtwixt.ConnectionString{
				Host: []string{"1"},
				Port: []string{"1", "2"},
			})
			assert.Error(t, err)
		})

		t.Run("HostAddr", func(t *testing.T) {
			_, err := Connector{}.Dialers(pgtwixt.ConnectionString{
				HostAddr: []string{"1"},
				Port:     []string{"1", "2"},
			})
			assert.Error(t, err)
		})
	})

	t.Run("TooFewPort", func(t *testing.T) {
		t.Run("Host", func(t *testing.T) {
			_, err := Connector{}.Dialers(pgtwixt.ConnectionString{
				Host: []string{"1", "2", "3"},
				Port: []string{"1", "2"},
			})
			assert.Error(t, err)
		})

		t.Run("HostAddr", func(t *testing.T) {
			_, err := Connector{}.Dialers(pgtwixt.ConnectionString{
				HostAddr: []string{"1", "2", "3"},
				Port:     []string{"1", "2"},
			})
			assert.Error(t, err)
		})
	})
}

func TestDialersMultipleHost(t *testing.T) {
	t.Parallel()

	t.Run("DefaultPort", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			Host: []string{"example.com", "/var/run/postgresql"},
		})
		require.NoError(t, err)
		require.Len(t, ds, 2)

		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "example.com:5432",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
				ServerName:    "example.com",
			},
		}, ds[0])
		assert.Equal(t, pgtwixt.UnixDialer{
			Address: "/var/run/postgresql/.s.PGSQL.5432",
		}, ds[1])
	})

	t.Run("OnePort", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			Host: []string{"example.com", "/var/run/postgresql"},
			Port: []string{"1000"},
		})
		require.NoError(t, err)
		require.Len(t, ds, 2)

		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "example.com:1000",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
				ServerName:    "example.com",
			},
		}, ds[0])
		assert.Equal(t, pgtwixt.UnixDialer{
			Address: "/var/run/postgresql/.s.PGSQL.1000",
		}, ds[1])
	})

	t.Run("TwoPort", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			Host: []string{"example.com", "/var/run/postgresql"},
			Port: []string{"1000", "2000"},
		})
		require.NoError(t, err)
		require.Len(t, ds, 2)

		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "example.com:1000",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
				ServerName:    "example.com",
			},
		}, ds[0])
		assert.Equal(t, pgtwixt.UnixDialer{
			Address: "/var/run/postgresql/.s.PGSQL.2000",
		}, ds[1])
	})
}

func TestDialersMultipleHostAddr(t *testing.T) {
	t.Parallel()

	t.Run("DefaultPort", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			HostAddr: []string{"127.0.0.1", "::1"},
		})
		require.NoError(t, err)
		require.Len(t, ds, 2)

		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "127.0.0.1:5432",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
			},
		}, ds[0])
		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "[::1]:5432",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
			},
		}, ds[1])
	})

	t.Run("OnePort", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			HostAddr: []string{"127.0.0.1", "::1"},
			Port:     []string{"1000"},
		})
		require.NoError(t, err)
		require.Len(t, ds, 2)

		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "127.0.0.1:1000",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
			},
		}, ds[0])
		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "[::1]:1000",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
			},
		}, ds[1])
	})

	t.Run("TwoPort", func(t *testing.T) {
		ds, err := Connector{}.Dialers(pgtwixt.ConnectionString{
			HostAddr: []string{"127.0.0.1", "::1"},
			Port:     []string{"1000", "2000"},
		})
		require.NoError(t, err)
		require.Len(t, ds, 2)

		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "127.0.0.1:1000",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
			},
		}, ds[0])
		assert.Equal(t, pgtwixt.TCPDialer{
			Address: "[::1]:2000",
			SSLConfig: tls.Config{
				MinVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateFreelyAsClient,
			},
		}, ds[1])
	})
}

package pgtwixt

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uhoh-itsmaciek/femebe/core"
)

const localhostServerCert = `-----BEGIN CERTIFICATE-----
MIIFUzCCAzugAwIBAgIJAIXWofFh1ke4MA0GCSqGSIb3DQEBCwUAMEAxCzAJBgNV
BAYTAlVTMQswCQYDVQQIDAJVUzEQMA4GA1UECgwHQ29tcGFueTESMBAGA1UEAwwJ
bG9jYWxob3N0MB4XDTE4MDEyMzAzNTEyMVoXDTI4MDEyMTAzNTEyMVowQDELMAkG
A1UEBhMCVVMxCzAJBgNVBAgMAlVTMRAwDgYDVQQKDAdDb21wYW55MRIwEAYDVQQD
DAlsb2NhbGhvc3QwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCWD8no
qiHcV0PjDkGLdfixFPDQwoQChiAic+MSTIL+nZmwGQPoMQg4Qo39j0kSmlxG+j4E
Sq6ZGmOgITKyv3rcb8OTUXuOA6wBLzyPWaUoRdIar9FJ7eBQdw+7SLKmpnBYy5lo
SrG2jeQPhVj6jkUFnNaXMzmuNxd7Vhg2VlvVvzfWeHrPF/xcpki4XrjzLlyjtCNt
h8hSfnSv3MZ0dA6TxEL/Ma5134e+xTj+V06P/czsmBnCqV1YPIS8NDErPh9hQXOj
YnI0a8/hH9gnYIUWQ61QzBNDpI/xGAJlw5idB3AvosY8uzTqdTDlScTy5VdCKCw7
wk+C+RB4A2oE7dOgN1rJ46LHJl6UlaYTQToOyNAbFHhKUohU4b869PXfb3xuQSve
xpABxaAMG4GZbBrpLNo4Nef4+a/zYkxoCtnusSaB89Q1fid0/czN93KgK455xUnB
Y+uzAu43E24n2cyT7sveqQz1iCYNoSaMDUhoyykXo8aD++KKc2hja9tx0u1GN8lI
0eYW4fxNpFZUqT6lYeTGWYJqWk5psPAiffguBM09H9+wF7qAtHHP5PyLr5x/SHQg
Btksw+Fp3B5WHoXaGL9va7hnmVQchgwDr++dVy9XS8IAJ8R1m/SRYzx+TeBkjVFY
L2NapBioC5SskwNq5p8s3hZSmGZpKYBAhoI9fQIDAQABo1AwTjAdBgNVHQ4EFgQU
AV1aYjFZjOKoUN19M6UHkUzN1CIwHwYDVR0jBBgwFoAUAV1aYjFZjOKoUN19M6UH
kUzN1CIwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEADJwtamIwSU4Q
37VM0hXvWgStYG2rN5X0HjWe6P546Z/0I3dtd5EDy65dgT3sxoLGU8dKLshcwkb4
OMYLQjUvwDR2bMnlifLpHvjVdKepPmvUAPuSwleY0gDruBg/t1IHlopbP6Lkji5k
O0SvPpTThmDAmRwE3dR5Kz0/RJJb5zNWGo9oZQ1lFmj6/eVnVvS3Z6qXV4MyAxtl
Xw+Y30DSq7A5NURZM2XbnqlKCwLRrPbtk4utcZ2RWNPXsW076JCQylHqT7xHbyLL
LxUnGwITepQMaBK11Y/WPbBQT726LmtIKnLTyhwcraGHm2GhEBZxQulojIG99ZTA
CcY7scEuNKuHq2O/QoxwXVpi2teXA3BezVxN3dXzGPwbHBpniHIvx25l4HMsk3CW
2qHoqSxCUvX6w+I7hMsJEUZ2LEpVFfJ3MPoZ1eUrhdaIbjRQw6M5FUGRKnjfRlUW
Ir9UFeIXBe/+G7+HR1ly88aRabhja/ykpTMHaz5oCwa7eTAqT5Sxn1sziVzCttZb
Z/YlK/rq52tKjbd490N13/LvR5yu+dY0SakMn+DpxKxcDdC8A7su+4E/739WcHyV
ETd8je71PqoJ6Qb9XNYZ16grCh2W4XDYfFuuASoohzaJNSyshPRfxGBSqMl5UuqU
RjEc+kCeqQuhMo9X1D7BQ8qpWS1oGTU=
-----END CERTIFICATE-----`

const localhostServerKey = `-----BEGIN PRIVATE KEY-----
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
ZImK7BgsSTBAc60aNwn8O99/cAopSQ==
-----END PRIVATE KEY-----`

func TestTCPDialerVerify(t *testing.T) {
	t.Parallel()

	listener, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer listener.Close()

	pool, err := x509.SystemCertPool()
	require.NoError(t, err)
	require.True(t, pool.AppendCertsFromPEM([]byte(localhostServerCert)))

	respond := func(response byte) {
		conn, err := listener.Accept()
		require.NoError(t, err)
		defer conn.Close()

		sslRequest := make([]byte, 8)
		n, err := io.ReadFull(conn, sslRequest)
		require.NoError(t, err)
		require.Equal(t, 8, n)

		n, err = conn.Write([]byte{response})
		require.NoError(t, err)
		require.Equal(t, 1, n)

		if response == core.AcceptSSLRequest {
			tls.Server(conn, &tls.Config{
				GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
					cert, err := tls.X509KeyPair([]byte(localhostServerCert), []byte(localhostServerKey))
					return &cert, err
				},
			}).Handshake()
		}
	}

	t.Run("Require", func(t *testing.T) {
		t.Run("Right", func(t *testing.T) {
			d := TCPDialer{
				Addr:    listener.Addr().String(),
				SSLMode: "require",
				SSLConfig: tls.Config{
					InsecureSkipVerify: true,
				},
			}

			go respond(core.AcceptSSLRequest)
			_, err := d.Dial(context.Background())
			assert.NoError(t, err)
		})

		t.Run("Wrong", func(t *testing.T) {
			d := TCPDialer{
				Addr:    listener.Addr().String(),
				SSLMode: "require",
				SSLConfig: tls.Config{
					InsecureSkipVerify: true,
				},
			}

			go respond(core.RejectSSLRequest)
			_, err := d.Dial(context.Background())
			assert.Error(t, err)
		})
	})

	t.Run("CA", func(t *testing.T) {
		t.Run("Right", func(t *testing.T) {
			d := TCPDialer{
				Addr:    listener.Addr().String(),
				SSLMode: "verify-ca",
				SSLConfig: tls.Config{
					InsecureSkipVerify: true,
					RootCAs:            pool,
				},
			}

			go respond(core.AcceptSSLRequest)
			_, err := d.Dial(context.Background())
			assert.NoError(t, err)
		})

		t.Run("None", func(t *testing.T) {
			d := TCPDialer{
				Addr:    listener.Addr().String(),
				SSLMode: "verify-ca",
				SSLConfig: tls.Config{
					InsecureSkipVerify: true,
					RootCAs:            pool,
				},
			}

			go respond(core.RejectSSLRequest)
			_, err := d.Dial(context.Background())
			assert.Error(t, err)
		})

		t.Run("Wrong", func(t *testing.T) {
			d := TCPDialer{
				Addr:    listener.Addr().String(),
				SSLMode: "verify-ca",
				SSLConfig: tls.Config{
					InsecureSkipVerify: true,
				},
			}

			go respond(core.AcceptSSLRequest)
			_, err := d.Dial(context.Background())
			assert.Error(t, err)
		})
	})

	t.Run("Full", func(t *testing.T) {
		t.Run("Right", func(t *testing.T) {
			d := TCPDialer{
				Addr:    listener.Addr().String(),
				SSLMode: "verify-full",
				SSLConfig: tls.Config{
					RootCAs:    pool,
					ServerName: "localhost",
				},
			}

			go respond(core.AcceptSSLRequest)
			_, err := d.Dial(context.Background())
			assert.NoError(t, err)
		})

		t.Run("None", func(t *testing.T) {
			d := TCPDialer{
				Addr:    listener.Addr().String(),
				SSLMode: "verify-full",
				SSLConfig: tls.Config{
					RootCAs:    pool,
					ServerName: "localhost",
				},
			}

			go respond(core.RejectSSLRequest)
			_, err := d.Dial(context.Background())
			assert.Error(t, err)
		})

		t.Run("UntrustedCert", func(t *testing.T) {
			d := TCPDialer{
				Addr:    listener.Addr().String(),
				SSLMode: "verify-full",
				SSLConfig: tls.Config{
					ServerName: "localhost",
				},
			}

			go respond(core.AcceptSSLRequest)
			_, err := d.Dial(context.Background())
			assert.Error(t, err)
		})

		t.Run("WrongHostname", func(t *testing.T) {
			d := TCPDialer{
				Addr:    listener.Addr().String(),
				SSLMode: "verify-full",
				SSLConfig: tls.Config{
					RootCAs:    pool,
					ServerName: "nope",
				},
			}

			go respond(core.AcceptSSLRequest)
			_, err := d.Dial(context.Background())
			assert.Error(t, err)
		})
	})
}

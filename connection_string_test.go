package pgtwixt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectionString(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		input    string
		expected ConnectionString
	}{
		{`host=a hostaddr=b port=c`, ConnectionString{
			Host: []string{"a"}, HostAddr: []string{"b"}, Port: []string{"c"},
		}},
		{`host=a,b,c hostaddr=d,e,f port=g,h,i`, ConnectionString{
			Host:     []string{"a", "b", "c"},
			HostAddr: []string{"d", "e", "f"},
			Port:     []string{"g", "h", "i"},
		}},
		{`dbname=x`, ConnectionString{Database: "x"}},
		{`user=jk password=lm`, ConnectionString{User: "jk", Password: "lm"}},
		{`passfile=/some/file`, ConnectionString{PasswordPath: "/some/file"}},
		{`connect_timeout = 5`, ConnectionString{ConnectTimeout: "5"}},
		{`client_encoding=Etc/UTC`, ConnectionString{ClientEncoding: "Etc/UTC"}},
		{`options='-c geqo=off'`, ConnectionString{Options: "-c geqo=off"}},
		{`application_name=foo`, ConnectionString{ApplicationName: "foo"}},
		{`fallback_application_name=bar`, ConnectionString{FallbackName: "bar"}},
		{`keepalives=0`, ConnectionString{KeepAlives: "0"}},
		{`keepalives_idle=1`, ConnectionString{KeepAlivesIdle: "1"}},
		{`keepalives_interval=2`, ConnectionString{KeepAlivesInterval: "2"}},
		{`keepalives_count=3`, ConnectionString{KeepAlivesCount: "3"}},
		{`sslmode=require`, ConnectionString{SSLMode: "require"}},
		{`requiressl=1`, ConnectionString{SSLRequire: "1"}},
		{`sslcompression=0`, ConnectionString{SSLCompression: "0"}},
		{`sslcert=some.crt`, ConnectionString{SSLCertPath: "some.crt"}},
		{`sslkey=some.key`, ConnectionString{SSLKeyPath: "some.key"}},
		{`sslrootcert=other.crt`, ConnectionString{SSLCAPath: "other.crt"}},
		{`sslcrl=other.crl`, ConnectionString{SSLCRLPath: "other.crl"}},
		{`requirepeer=postgres`, ConnectionString{RequirePeer: "postgres"}},
		{`service=baz`, ConnectionString{Service: "baz"}},
		{`unknown=val many=times`, ConnectionString{
			Remainder: map[string]string{"unknown": "val", "many": "times"},
		}},
	} {
		t.Run(tt.input, func(t *testing.T) {
			if tt.expected.Remainder == nil {
				tt.expected.Remainder = make(map[string]string)
			}

			var result ConnectionString
			require.NoError(t, result.Parse(tt.input))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConnectionStringBadQuoting(t *testing.T) {
	t.Parallel()

	var result ConnectionString
	assert.Error(t, result.Parse(`a`))
	assert.Error(t, result.Parse(`a=`))
	assert.Error(t, result.Parse(`a='`))
}

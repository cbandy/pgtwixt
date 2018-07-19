package pgtwixt

import (
	"context"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnixDialerInterface(t *testing.T) {
	var _ Dialer = UnixDialer{}
}

func TestUnixDialerRequirePeer(t *testing.T) {
	t.Parallel()

	dir, err := ioutil.TempDir("", "pgwtixt-unix-dialer")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	spath := path.Join(dir, "server")
	listener, err := net.Listen("unix", spath)
	require.NoError(t, err)
	defer listener.Close()

	t.Run("None", func(t *testing.T) {
		d := UnixDialer{Address: spath}
		_, err := d.Dial(context.Background())
		assert.NoError(t, err)
	})

	u, err := user.Current()
	require.NoError(t, err)

	t.Run("Right", func(t *testing.T) {
		d := UnixDialer{Address: spath, RequirePeer: u.Username}
		_, err := d.Dial(context.Background())
		assert.NoError(t, err)
	})

	t.Run("Wrong", func(t *testing.T) {
		d := UnixDialer{Address: spath, RequirePeer: "nope" + u.Username}
		_, err := d.Dial(context.Background())
		assert.Error(t, err)
	})
}

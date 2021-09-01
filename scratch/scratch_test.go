// +build !nautilus

// Initially, we're only providing mirroring related functions for octopus as
// that version of ceph deprecated a number of the functions in nautilus. If
// you need mirroring on an earlier supported version of ceph please file an
// issue in our tracker.

package scratch

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ceph/go-ceph/rados"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/ceph/go-ceph/rbd"
)

func radosConnect(t *testing.T) *rados.Conn {
	conn, err := rados.NewConn()
	require.NoError(t, err)
	err = conn.ReadDefaultConfigFile()
	require.NoError(t, err)
	waitForRadosConn(t, conn)
	return conn
}

func radosConnectConfig(t *testing.T, p string) *rados.Conn {
	conn, err := rados.NewConn()
	require.NoError(t, err)
	err = conn.ReadConfigFile(p)
	require.NoError(t, err)
	waitForRadosConn(t, conn)
	return conn
}

func waitForRadosConn(t *testing.T, conn *rados.Conn) {
	var err error
	timeout := time.After(time.Second * 15)
	ch := make(chan error)
	go func(conn *rados.Conn) {
		ch <- conn.Connect()
	}(conn)
	select {
	case err = <-ch:
	case <-timeout:
		err = fmt.Errorf("timed out waiting for connect")
	}
	require.NoError(t, err)
}

func TestGetMirrorUUID(t *testing.T) {
	conn := radosConnect(t)
	poolName := "a"
	err := conn.MakePool(poolName)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, conn.DeletePool(poolName))
		conn.Shutdown()
	}()

	for i := range []int{1, 2, 3, 5, 0, 0, 0, 0} {

		name1 := "img" + strconv.Itoa(i)
		go func(name1 string) {
			conn := radosConnect(t)
			ioctx, err := conn.OpenIOContext(poolName)
			assert.NoError(t, err)
			defer func() {
				ioctx.Destroy()
			}()
			fmt.Printf("creating %s\n", name1)
			options := NewRbdImageOptions()
			err = CreateImage(ioctx, name1, uint64(100), options)
			require.NoError(t, err)

			fmt.Printf("done creating %s\n", name1)
		}(name1)
	}
	time.Sleep(time.Second * 30)
	panic("hello")
}

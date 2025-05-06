package cephfs

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	fsadmin "github.com/ceph/go-ceph/cephfs/admin"
	"github.com/ceph/go-ceph/internal/admintest"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

var (
	radosConnector = admintest.NewConnector()
)

// NoGroup should be used when an optional subvolume group name is not
// specified.
const NoGroup = ""

func TestDiff(t *testing.T) {
	fsa := fsadmin.NewFromConn(radosConnector.Get(t))
	volume := "cephfs"

	subname := "SubVol1"
	err := fsa.CreateSubVolume(volume, NoGroup, subname, nil)
	assert.NoError(t, err)
	defer func() {
		err := fsa.RemoveSubVolume(volume, NoGroup, subname)
		assert.NoError(t, err)
	}()

	path, err := fsa.SubVolumePath(volume, NoGroup, subname)
	assert.NoError(t, err)
	t.Logf("path: %v", path)

	mount := fsConnect(t)
	defer fsDisconnect(t, mount)

	t.Log("getting client_snapdir")
	t.Log(mount.GetConfigOption("client_snapdir"))

	t.Log("setting client_snapdir", mount.SetConfigOption("client_snapdir", ".snap"))

	t.Log("getting client_snapdir")
	t.Log(mount.GetConfigOption("client_snapdir"))

	t.Log("getting debug_client=20")
	t.Log(mount.GetConfigOption("debug_client"))

	t.Log("setting debug_client=20", mount.SetConfigOption("debug_client", "20"))
	t.Log("getting debug_client=20")
	t.Log(mount.GetConfigOption("debug_client"))

	t.Log("getting log_file")
	t.Log(mount.GetConfigOption("log_file"))
	t.Log("setting log_file", mount.SetConfigOption("log_file", "/tmp/cephfs.log"))
	t.Log("getting log_file")
	t.Log(mount.GetConfigOption("log_file"))
	// mount, err := CreateMount()
	// require.NoError(t, err)
	// require.NotNil(t, mount)
	// defer func() {
	// 	assert.NoError(t, mount.Release())
	// }()

	// err = mount.ReadDefaultConfigFile()
	// require.NoError(t, err)

	// err = mount.MountWithRoot("/")
	// assert.NoError(t, err)
	// defer func() {
	// 	assert.NoError(t, mount.Unmount())
	// }()

	err = DoFileOp(mount, path, 10, 20, 0, OpCreate)
	assert.NoError(t, err)

	snap1 := "Snap1"
	err = fsa.CreateSubVolumeSnapshot(volume, NoGroup, subname, snap1)
	assert.NoError(t, err)
	defer func() {
		err := fsa.RemoveSubVolumeSnapshot(volume, NoGroup, subname, snap1)
		assert.NoError(t, err)
	}()

	err = DoFileOp(mount, path, 10, 20, 20, OpModify)
	assert.NoError(t, err)

	snap2 := "Snap2"
	err = fsa.CreateSubVolumeSnapshot(volume, NoGroup, subname, snap2)
	assert.NoError(t, err)
	defer func() {
		err := fsa.RemoveSubVolumeSnapshot(volume, NoGroup, subname, snap2)
		assert.NoError(t, err)
	}()

	// t.Log(mount.CurrentDir())
	// dirPaths := []string{"/volumes"}
	// newDirPaths := []string{}
	// for {
	// 	for _, dirPath := range dirPaths {
	// 		t.Logf("dirPath: %v", dirPath)
	// 		Dir, err := mount.OpenDir(dirPath)
	// 		if err != nil {
	// 			t.Logf("open dir %v: %v", dirPath, err)
	// 			continue
	// 		}
	// 		for {
	// 			dirEntry, err := Dir.ReadDir()
	// 			if err != nil {
	// 				t.Log(err)
	// 				continue
	// 			}
	// 			if dirEntry == nil {
	// 				break
	// 			}
	// 			if dirEntry.Name() == "." || dirEntry.Name() == ".." {
	// 				continue
	// 			}
	// 			t.Logf("dirEntry: %v: %v: %v", dirEntry.Name(), dirEntry.Inode(), dirEntry.DType())
	// 			if dirEntry.DType() == DTypeDir {
	// 				newDirPaths = append(newDirPaths, dirPath+"/"+dirEntry.Name())
	// 			}
	// 		}
	// 	}
	// 	if len(newDirPaths) == 0 {
	// 		break
	// 	}
	// 	dirPaths = newDirPaths
	// 	newDirPaths = []string{}
	// }
	// dirPaths = []string{"/volumes/_nogroup/SubVol1/.snap"}
	// newDirPaths = []string{}
	// for {
	// 	for _, dirPath := range dirPaths {
	// 		t.Logf("dirPath: %v", dirPath)
	// 		Dir, err := mount.OpenDir(dirPath)
	// 		if err != nil {
	// 			t.Logf("open dir %v: %v", dirPath, err)
	// 			continue
	// 		}
	// 		for {
	// 			dirEntry, err := Dir.ReadDir()
	// 			if err != nil {
	// 				t.Log(err)
	// 				continue
	// 			}
	// 			if dirEntry == nil {
	// 				break
	// 			}
	// 			if dirEntry.Name() == "." || dirEntry.Name() == ".." {
	// 				continue
	// 			}
	// 			t.Logf("dirEntry: %v: %v: %v", dirEntry.Name(), dirEntry.Inode(), dirEntry.DType())
	// 			if dirEntry.DType() == DTypeDir {
	// 				newDirPaths = append(newDirPaths, dirPath+"/"+dirEntry.Name())
	// 			}
	// 		}
	// 	}
	// 	if len(newDirPaths) == 0 {
	// 		break
	// 	}
	// 	dirPaths = newDirPaths
	// 	newDirPaths = []string{}
	// }

	snap1ID, err := GetSnapshotID(mount, "/volumes/_nogroup/SubVol1/.snap/"+snap1)
	assert.NoError(t, err)
	snap2ID, err := GetSnapshotID(mount, "/volumes/_nogroup/SubVol1/.snap/"+snap2)
	assert.NoError(t, err)
	t.Logf("snap1ID: %v", snap1ID)
	t.Logf("snap2ID: %v", snap2ID)
	err = mount.ChangeDir("/")
	assert.NoError(t, err)
	t.Log(mount.CurrentDir())
	t.Log(path)
	// err = mount.ChangeDir(path)
	// assert.NoError(t, err)

	t.Log(mount.CurrentDir())
	t.Logf("rootPath: %q", "/volumes/_nogroup/SubVol1/")
	t.Logf("relPath: %q", "/")
	diff, err := CephOpenSnapDiff(SnapDiffConfig{
		CMount:   mount,
		RootPath: "/volumes/_nogroup/SubVol1/",
		RelPath:  "/",
		Snap1:    snap1,
		Snap2:    snap2,
	})
	assert.NoError(t, err)
	assert.NotNil(t, diff)

	defer func() {
		err = CephCloseSnapDiff(diff)
		t.Logf("close snap diff: %v", err)
	}()

	t.Logf("diff: %v", diff)
	assert.NotNil(t, diff.Dir1)
	assert.NotNil(t, diff.DirAux)
	t.Logf("mount: %q", diff.CMount.CurrentDir())

	for {
		diffEntry, err := CephReaddirSnapDiff(diff)
		if err != nil {
			t.Logf("readdir snap diff error: %v", err)
			break
		}
		if diffEntry == nil {
			break
		}
		t.Logf("diffEntry: %v: %v: %v", diffEntry.DirEntry.Name(), diffEntry.DirEntry.Inode(),
			diffEntry.DirEntry.DType())
		if diffEntry.DirEntry.Name() == "." || diffEntry.DirEntry.Name() == ".." {
			continue
		}
	}
}

type FileOp string

const (
	OpCreate FileOp = "create"
	OpDelete FileOp = "delete"
	OpModify FileOp = "modify"
)

func DoFileOp(mount *MountInfo, path string, numFiles int, numBytes int, percent int, Op FileOp) error {
	for i := 0; i < numFiles; i++ {
		if Op == OpCreate || rand.Intn(100) >= percent {
			break
		}

		fileName := path + "/file-" + strconv.Itoa(i) + ".txt"
		// Create the file
		f, err := mount.Open(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to create file %s: %w", fileName, err)
		}
		defer func() {
			f.Sync()
			f.Close()
		}()
		buf := make([]byte, numBytes)
		_, err = rand.Read(buf)
		if err != nil {
			return fmt.Errorf("failed to read random bytes: %w", err)
		}

		_, err = f.write(buf, 0)
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %w", fileName, err)
		}
	}
	return nil
}

/*
=== RUN   TestDiff
    diff_test.go:50: /
    diff_test.go:55: dirPath: /volumes
    diff_test.go:73: dirEntry: _:SubVol1.meta: 1099511628283: 8
    diff_test.go:73: dirEntry: _nogroup: 1099511627779: 4
    diff_test.go:55: dirPath: /volumes/_nogroup
    diff_test.go:73: dirEntry: SubVol1: 1099511627780: 4
    diff_test.go:55: dirPath: /volumes/_nogroup/SubVol1
    diff_test.go:73: dirEntry: 5657dd89-9a75-470d-84e1-3483f0e220f1: 1099511627781: 4
    diff_test.go:73: dirEntry: .meta: 1099511627782: 8
    diff_test.go:55: dirPath: /volumes/_nogroup/SubVol1/5657dd89-9a75-470d-84e1-3483f0e220f1
    diff_test.go:89: dirPath: /volumes/_nogroup/SubVol1/.snap
    diff_test.go:107: dirEntry: Snap1: 1099511627780: 4
    diff_test.go:107: dirEntry: Snap2: 1099511627780: 4
    diff_test.go:89: dirPath: /volumes/_nogroup/SubVol1/.snap/Snap1
    diff_test.go:107: dirEntry: 5657dd89-9a75-470d-84e1-3483f0e220f1: 1099511627781: 4
    diff_test.go:107: dirEntry: .meta: 1099511627782: 8
    diff_test.go:89: dirPath: /volumes/_nogroup/SubVol1/.snap/Snap2
    diff_test.go:107: dirEntry: 5657dd89-9a75-470d-84e1-3483f0e220f1: 1099511627781: 4
    diff_test.go:107: dirEntry: .meta: 1099511627782: 8
    diff_test.go:89: dirPath: /volumes/_nogroup/SubVol1/.snap/Snap1/5657dd89-9a75-470d-84e1-3483f0e220f1
    diff_test.go:89: dirPath: /volumes/_nogroup/SubVol1/.snap/Snap2/5657dd89-9a75-470d-84e1-3483f0e220f1
*/

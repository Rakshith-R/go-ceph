package rbd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetSnapshotByID(t *testing.T) {
	conn := radosConnect(t)

	poolname := GetUUID()
	err := conn.MakePool(poolname)
	assert.NoError(t, err)

	ioctx, err := conn.OpenIOContext(poolname)
	require.NoError(t, err)

	name := GetUUID()
	snapName := fmt.Sprintf("snap-%s", GetUUID())
	cloneName := fmt.Sprintf("clone-%s", GetUUID())
	groupName := fmt.Sprintf("group-%s", GetUUID())
	options := NewRbdImageOptions()
	assert.NoError(t,
		options.SetUint64(ImageOptionOrder, uint64(testImageOrder)))
	err = CreateImage(ioctx, name, testImageSize, options)
	assert.NoError(t, err)

	img, err := OpenImage(ioctx, name, NoSnapshot)
	assert.NoError(t, err)

	t.Run("set non-existent snapshot by id", func(t *testing.T) {
		err = img.SetSnapshotByID(0)
		assert.Error(t, err)
	})

	t.Run("set regular snapshot by id", func(t *testing.T) {
		snap, err := img.CreateSnapshot("snap")
		assert.NoError(t, err)
		defer func() { assert.NoError(t, snap.Remove()) }()

		snapID, err := img.GetSnapID(snap.name)
		assert.NoError(t, err)

		err = img.SetSnapshotByID(snapID)
		assert.NoError(t, err)
	})

	t.Run("set regular snapshot in trash by id", func(t *testing.T) {
		snap, err := img.CreateSnapshot("snap")
		assert.NoError(t, err)
		defer func() { assert.NoError(t, snap.Remove()) }()

		snapID, err := img.GetSnapID(snap.name)
		assert.NoError(t, err)

		err = CloneImage(ioctx, img.name, snap.name, ioctx, "clone", options)
		assert.NoError(t, err)
		defer func() { assert.NoError(t, RemoveImage(ioctx, "clone")) }()

		assert.NoError(t, snap.Remove())

		err = img.SetSnapshotByID(snapID)
		assert.NoError(t, err)
	})

	t.Run("set group snapshot by id", func(t *testing.T) {
		err = GroupCreate(ioctx, groupName)
		assert.NoError(t, err)
		defer func() { assert.NoError(t, GroupRemove(ioctx, "group")) }()

		err = GroupImageAdd(ioctx, groupName, ioctx, name)
		assert.NoError(t, err)
		defer func() { assert.NoError(t, GroupImageRemove(ioctx, groupName, ioctx, name)) }()

		err = GroupSnapCreate(ioctx, groupName, snapName)
		assert.NoError(t, err)
		defer func() { assert.NoError(t, GroupSnapRemove(ioctx, groupName, snapName)) }()

		groupSnapInfo, err := GroupSnapGetInfo(ioctx, groupName, snapName)
		assert.NoError(t, err)

		snapID := groupSnapInfo.Snapshots[0].SnapID

		err = img.SetSnapshotByID(snapID)
		assert.NoError(t, err)
	})

	t.Run("set group snapshot in trash by id", func(t *testing.T) {
		err = GroupCreate(ioctx, groupName)
		assert.NoError(t, err)
		defer func() { assert.NoError(t, GroupRemove(ioctx, "group")) }()

		err = GroupImageAdd(ioctx, groupName, ioctx, name)
		assert.NoError(t, err)
		defer func() { assert.NoError(t, GroupImageRemove(ioctx, groupName, ioctx, name)) }()

		err = GroupSnapCreate(ioctx, groupName, snapName)
		assert.NoError(t, err)

		groupSnapInfo, err := GroupSnapGetInfo(ioctx, groupName, snapName)
		assert.NoError(t, err)

		snapID := groupSnapInfo.Snapshots[0].SnapID

		err = CloneImageByID(ioctx, img.name, snapID, ioctx, cloneName, options)
		assert.NoError(t, err)
		defer func() { assert.NoError(t, RemoveImage(ioctx, cloneName)) }()

		err = GroupSnapRemove(ioctx, groupName, snapName)
		assert.NoError(t, err)

		err = img.SetSnapshotByID(snapID)
		assert.NoError(t, err)
	})
}

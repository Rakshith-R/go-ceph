// +build !nautilus

package admin

import (
	"fmt"

	ccom "github.com/ceph/go-ceph/common/commands"
)

// RBDAdmin is used to administrate rbd volumes and pools.
type RBDAdmin struct {
	conn ccom.RadosCommander
}

// NewFromConn creates an new management object from a preexisting
// rados connection. The existing connection can be rados.Conn or any
// type implementing the RadosCommander interface.
func NewFromConn(conn ccom.RadosCommander) *RBDAdmin {
	return &RBDAdmin{conn}
}

// LevelSpec values are used to identify RBD objects wherever Ceph APIs
// require a levelspec to select an image, pool, or namespace.
type LevelSpec struct {
	spec string
}

// NewLevelSpec is used to construct a LevelSpec given a pool and
// optional namespace and image names.
func NewLevelSpec(pool, namespace, image string) LevelSpec {
	var s string
	if image != "" && namespace != "" {
		s = fmt.Sprintf("%s/%s/%s", pool, namespace, image)
	} else if image != "" {
		s = fmt.Sprintf("%s/%s", pool, image)
	} else if namespace != "" {
		s = fmt.Sprintf("%s/%s/", pool, namespace)
	} else {
		s = fmt.Sprintf("%s/", pool)
	}
	return LevelSpec{s}
}

// NewRawLevelSpec returns a LevelSpec directly based on the spec string
// argument without constructing it from component values. This should only be
// used if NewLevelSpec can not create the levelspec value you want to pass to
// ceph.
func NewRawLevelSpec(spec string) LevelSpec {
	return LevelSpec{spec}
}

// parseImageSpec is used to construct a ImageSpec given a image and
// optional namespace and pool names.
func parseImageSpec(pool, namespace, image string) string {
	var s string
	if pool != "" && namespace != "" {
		s = fmt.Sprintf("%s/%s/%s", pool, namespace, image)
	} else if pool != "" {
		s = fmt.Sprintf("%s/%s", pool, image)
	} else {
		s = image
	}
	return s
}

// ImageSpec values are used to identify RBD objects wherever Ceph APIs
// require a imagespec to select an image, pool, or namespace.
type ImageSpec struct {
	spec string
}

// NewImageSpec is used to construct a ImageSpec given a image and
// optional namespace and pool names.
func NewImageSpec(pool, namespace, image string) ImageSpec {
	return ImageSpec{parseImageSpec(pool, namespace, image)}
}

// NewRawImageSpec returns a ImageSpec directly based on the spec string
// argument without constructing it from component values. This should only be
// used if NewImageSpec can not create the imagespec value you want to pass to
// ceph.
func NewRawImageSpec(spec string) ImageSpec {
	return ImageSpec{spec}
}

// ImageIdSpecvalues are used to identify RBD objects wherever Ceph APIs
// require a imagespec to select an image, pool, or namespace.
type ImageIdSpec struct {
	spec string
}

// NewImageIdSpec is used to construct a ImageIdSpec given a imageId and
// optional namespace and pool names.
func NewImageIdSpec(pool, namespace, imageId string) ImageIdSpec {
	return ImageIdSpec{parseImageSpec(pool, namespace, imageId)}
}

// NewRawImageIdSpec returns a ImageIdSpec directly based on the spec string
// argument without constructing it from component values. This should only be
// used if NewImageIdSpec can not create the imagespec value you want to pass to
// ceph.
func NewRawImageIdSpec(spec string) ImageIdSpec {
	return ImageIdSpec{spec}
}

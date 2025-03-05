//go:build ceph_preview

package rbd

/*
#cgo LDFLAGS: -lrbd
#undef _GNU_SOURCE
#include <errno.h>
#include <stdlib.h>
#include <rbd/librbd.h>

extern int diffIterateByIDCallback(uint64_t, size_t, int, uintptr_t);

// rbd_diff_iterate3_fn matches the rbd_diff_iterate3 function signature.
typedef int(*rbd_diff_iterate3_fn)(rbd_image_t image, uint64_t from_snap_id,
	uint64_t ofs, uint64_t len, uint32_t flags,
	int (*cb)(uint64_t, size_t, int, void *), void *arg);

// rbd_diff_iterate3_dlsym take *fn as rbd_diff_iterate3_fn and calls the dynamically loaded
// rbd_diff_iterate3 function passed as 1st argument.
static inline int rbd_diff_iterate3_dlsym(void *fn, rbd_image_t image,
	uint64_t from_snap_id, uint64_t ofs, uint64_t len, uint32_t flags, uintptr_t arg) {
	// cast function pointer fn to rbd_diff_iterate3 and call the function
	return ((rbd_diff_iterate3_fn) fn)(image, from_snap_id, ofs, len, flags, (void*)diffIterateByIDCallback, (void*)arg);
}
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/ceph/go-ceph/internal/callbacks"
	"github.com/ceph/go-ceph/internal/dlsym"
)

const (
	rbdDiffIterateFlagIncludeParentKey = "RBD_DIFF_ITERATE_FLAG_INCLUDE_PARENT"
	rbdDiffIterateFlagWholeObjectKey   = "RBD_DIFF_ITERATE_FLAG_WHOLE_OBJECT"
)

var (
	diffIterateByIDCallbacks = callbacks.New()
	diffIterateByIDOnce      sync.Once
	diffIterateById          unsafe.Pointer
	diffIterateByIdErr       error
)

// DiffIterateByIDCallback defines the function signature needed for the
// DiffIterateByID callback.
//
// The function will be called with the arguments: offset, length, exists, and
// data. The offset and length correspond to the changed region of the image.
// The exists value is set to zero if the region is known to be zeros,
// otherwise it is set to 1. The data value is the extra data parameter that
// was set on the DiffIterateByIDConfig and is meant to be used for passing
// arbitrary user-defined items to the callback function.
//
// The callback can trigger the iteration to terminate early by returning
// a non-zero error code.
type DiffIterateByIDCallback func(uint64, uint64, int, interface{}) int

// DiffIterateByIDConfig is used to define the parameters of a DiffIterateByID call.
// Callback, Offset, and Length should always be specified when passed to
// DiffIterateByID. The other values are optional.
type DiffIterateByIDConfig struct {
	FromSnapID    uint64
	Offset        uint64
	Length        uint64
	IncludeParent DiffIncludeParent
	WholeObject   DiffWholeObject
	Callback      DiffIterateByIDCallback
	Data          interface{}
}

// DiffIterateByID calls a callback on changed extents of an image.
//
// Calling DiffIterateByID will cause the callback specified in the
// DiffIterateConfig to be called as many times as there are changed
// regions in the image (controlled by the parameters as passed to librbd).
//
// See the documentation of DiffIterateByIDCallback for a description of the
// arguments to the callback and the return behavior.
//
// Implements:
//
//	int rbd_diff_iterate3(rbd_image_t image,
//	                      uint64_t from_snap_id,
//	                      uint64_t ofs, uint64_t len,
//	                      uint32_t flags,
//	                      int (*cb)(uint64_t, size_t, int, void *),
//	                      void *arg);
func (image *Image) DiffIterateByID(config DiffIterateByIDConfig) error {
	if err := image.validate(imageIsOpen); err != nil {
		return err
	}
	if config.Callback == nil {
		return getError(C.EINVAL)
	}

	diffIterateByIDOnce.Do(func() {
		diffIterateById, diffIterateByIdErr = dlsym.LookupSymbol("rbd_diff_iterate3")
	})

	if diffIterateByIdErr != nil {
		return fmt.Errorf("%w: %w", ErrNotImplemented, diffIterateByIdErr)
	}

	cbIndex := diffIterateByIDCallbacks.Add(config)
	defer diffIterateByIDCallbacks.Remove(cbIndex)

	flags := C.uint32_t(0)
	if config.IncludeParent == IncludeParent {
		flags |= C.RBD_DIFF_ITERATE_FLAG_INCLUDE_PARENT
	}
	if config.WholeObject == EnableWholeObject {
		flags |= C.RBD_DIFF_ITERATE_FLAG_WHOLE_OBJECT
	}

	ret := C.rbd_diff_iterate3_dlsym(
		diffIterateById,
		image.image,
		C.uint64_t(config.FromSnapID),
		C.uint64_t(config.Offset),
		C.uint64_t(config.Length),
		flags,
		C.uintptr_t(cbIndex))

	return getError(ret)
}

//export diffIterateByIDCallback
func diffIterateByIDCallback(
	offset C.uint64_t, length C.size_t, exists C.int, index uintptr) C.int {

	v := diffIterateByIDCallbacks.Lookup(index)
	config := v.(DiffIterateByIDConfig)
	return C.int(config.Callback(
		uint64(offset), uint64(length), int(exists), config.Data))
}

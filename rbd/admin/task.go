// +build !nautilus

package admin

import (
	ccom "github.com/ceph/go-ceph/common/commands"
	"github.com/ceph/go-ceph/internal/commands"
)

// TaskAdmin encapsulates management functions for
// ceph rbd mirror snapshot schedules.
type TaskAdmin struct {
	conn ccom.MgrCommander
}

// TaskAdd returns a TaskAdmin type for
// managing ceph rbd task operations.
func (ra *RBDAdmin) Task() *TaskAdmin {
	return &TaskAdmin{conn: ra.conn}
}

/*
{
   "sequence":4,
   "id":"b3da71d2-410a-44dd-bf98-dc2d1b3238c2",
   "message":"Removing image replicapool/csi-snap-d586c4f8-1459-11ec-a6b2-0242ac110007",
   "refs":{
      "action":"remove",
      "pool_name":"replicapool",
      "pool_namespace":"",
      "image_name":"csi-snap-d586c4f8-1459-11ec-a6b2-0242ac110007",
      "image_id":"129b14807d52"
   },
   "in_progress":true,
   "progress":0.7171875,
   "retry_attempts":1,
   "retry_time":"2021-09-13T06:17:35.742165"
   "retry_message": "[errno 39] RBD image has snapshots (error deleting image from trash)",
}
*/
type TaskRefs struct {
	Action        string `json:"action"`
	PoolName      string `json:"pool_name"`
	PoolNamespace string `json:"pool_namespace"`
	ImageName     string `json:"image_name"`
	ImageID       string `json:"image_id"`
}

type TaskResponse struct {
	Sequence      int      `json:"sequence"`
	Id            string   `json:"id"`
	Message       string   `json:"message"`
	Refs          TaskRefs `json:"refs"`
	InProgress    bool     `json:"in_progress"`
	Progress      float64  `json:"progress"`
	RetryAttempts int      `json:"retry_attempts"`
	RetryTime     string   `json:"retry_time"`
	RetryMessage  string   `json:"retry_message"`
}

func parseTaskResponseList(res commands.Response) ([]TaskResponse, error) {
	var taskResponseList []TaskResponse
	err := res.NoStatus().Unmarshal(&taskResponseList).End()

	return taskResponseList, err
}

func parseTaskResponse(res commands.Response) (TaskResponse, error) {
	var taskResponse TaskResponse
	err := res.NoStatus().Unmarshal(&taskResponse).End()

	return taskResponse, err
}

// Add a task to flatten image based on give image_spec.
//
// Similar To:
//  rbd task add flatten <image_spec>
func (ta *TaskAdmin) AddFlatten(img ImageSpec) (TaskResponse, error) {
	m := map[string]string{
		"prefix":     "rbd task add flatten",
		"image_spec": img.spec,
		"format":     "json",
	}
	return parseTaskResponse(commands.MarshalMgrCommand(ta.conn, m))
}

// Add a task to remove image based on give image_spec.
//
// Similar To:
//  rbd task add remove <image_spec>
func (ta *TaskAdmin) AddRemove(img ImageSpec) (TaskResponse, error) {
	m := map[string]string{
		"prefix":     "rbd task add remove",
		"image_spec": img.spec,
		"format":     "json",
	}

	return parseTaskResponse(commands.MarshalMgrCommand(ta.conn, m))
}

// Add a task to trash remove image based on give image_id_spec.
//
// Similar To:
//  rbd task add trash remove <image_id_spec>
func (ta *TaskAdmin) AddTrashRemove(img ImageIdSpec) (TaskResponse, error) {
	m := map[string]string{
		"prefix":        "rbd task add trash remove",
		"image_id_spec": img.spec,
		"format":        "json",
	}

	return parseTaskResponse(commands.MarshalMgrCommand(ta.conn, m))
}

func (ta *TaskAdmin) AddMigrationCommit(img ImageSpec) (TaskResponse, error) {
	m := map[string]string{
		"prefix":     "rbd task add migration commit",
		"image_spec": img.spec,
		"format":     "json",
	}

	return parseTaskResponse(commands.MarshalMgrCommand(ta.conn, m))
}

func (ta *TaskAdmin) AddMigrationAbort(img ImageSpec) (TaskResponse, error) {
	m := map[string]string{
		"prefix":     "rbd task add migration abort",
		"image_spec": img.spec,
		"format":     "json",
	}

	return parseTaskResponse(commands.MarshalMgrCommand(ta.conn, m))
}

func (ta *TaskAdmin) AddMigrationExecute(img ImageSpec) (TaskResponse, error) {
	m := map[string]string{
		"prefix":     "rbd task add migration execute",
		"image_spec": img.spec,
		"format":     "json",
	}

	return parseTaskResponse(commands.MarshalMgrCommand(ta.conn, m))
}

func (ta *TaskAdmin) List() ([]TaskResponse, error) {
	m := map[string]string{
		"prefix": "rbd task list",
		"format": "json",
	}

	return parseTaskResponseList(commands.MarshalMgrCommand(ta.conn, m))
}

func (ta *TaskAdmin) GetTask(taskID string) (TaskResponse, error) {
	m := map[string]string{
		"prefix":  "rbd task list",
		"task_id": taskID,
		"format":  "json",
	}

	return parseTaskResponse(commands.MarshalMgrCommand(ta.conn, m))
}

func (ta *TaskAdmin) Cancel(taskID string) (TaskResponse, error) {
	m := map[string]string{
		"prefix":  "rbd task cancel",
		"task_id": taskID,
		"format":  "json",
	}

	return parseTaskResponse(commands.MarshalMgrCommand(ta.conn, m))
}

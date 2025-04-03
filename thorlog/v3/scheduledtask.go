package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type ScheduledTask struct {
	LogObjectHeader

	Name    string `json:"name" textlog:"name"`
	Path    string `json:"path" textlog:"path"`
	Command string `json:"command" textlog:"command"`
	Enabled bool   `json:"enabled" textlog:"enabled"`
	LastRun Time   `json:"lastrun,omitzero" textlog:"lastrun,omitempty"`
	NextRun Time   `json:"nextrun,omitzero" textlog:"nextrun,omitempty"`
}

func (ScheduledTask) reportable() {}

const typeScheduledTask = "scheduled task"

func init() { AddLogObjectType(typeScheduledTask, &ScheduledTask{}) }

func NewScheduledTask() *ScheduledTask {
	return &ScheduledTask{
		LogObjectHeader: LogObjectHeader{
			Type: typeScheduledTask,
		},
	}
}

type RegistryScheduledTask struct {
	jsonlog.ObjectHeader
	RegistryHive string `json:"registry_hive" textlog:"hive"`
	Key          string `json:"key" textlog:"registry_path"`
	Guid         string `json:"guid" textlog:"guid"`
	Path         string `json:"path" textlog:"path"`
	Version      int    `json:"version" textlog:"version"`
	Created      Time   `json:"created" textlog:"created"`
	LastRun      Time   `json:"last_run" textlog:"last_run"`
	LastStopped  Time   `json:"last_stopped" textlog:"last_stopped"`
	Status       string `json:"status" textlog:"status"`
	LastResult   string `json:"last_result" textlog:"last_result"`
}

func (RegistryScheduledTask) reportable() {}

const typeRegistryScheduledTask = "registry scheduled task"

func init() { AddLogObjectType(typeRegistryScheduledTask, &RegistryScheduledTask{}) }

func NewRegistryScheduledTask() *RegistryScheduledTask {
	return &RegistryScheduledTask{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeRegistryScheduledTask,
		},
	}
}

package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

// ScheduledTask describes a Windows Scheduled Task.
//
// See also the Microsoft documentation at https://learn.microsoft.com/en-us/windows/win32/taskschd/task-scheduler-reference
// for more details about scheduled tasks.
type ScheduledTask struct {
	LogObjectHeader

	// Name of the scheduled task.
	Name string `json:"name" textlog:"name"`
	// Path (within C:\Windows\System32\tasks) of this scheduled task.
	Path string `json:"path" textlog:"path"`

	// Commands executed when this scheduled task activates. Commands each include both image and arguments.
	Commands StringList `json:"commands" textlog:"command,omitempty"`
	// COM Handlers (as GUIDs) invoked when this scheduled task activates.
	ComHandlers StringList `json:"com_handlers,omitempty" textlog:"com_handler,expand,omitempty"`

	// Whether the scheduled task is active.
	Enabled bool `json:"enabled" textlog:"enabled"`
	// The trigger types when the task should be executed.
	// Options:
	// - Time (at a fixed time)
	// - Calendar (regularly based on calendar)
	// - Boot
	// - Logon
	// - Event (when specific events occur in the Windows Eventlog)
	// - Registration (only when the task was initially created)
	// - SessionStateChange (configurable on e.g. remote connection, session unlock, ...)
	Triggers StringList `json:"triggers,omitempty" textlog:"triggers,omitempty"`

	// The user (or SID) as which the scheduled task will run.
	User string `json:"user" textlog:"user"`
	// Logon type, options: S4U, Password, InteractiveToken
	LogonType string `json:"logon_type" textlog:"logon_type"`
	// Run level, options: LeastPrivilege or HighestAvailable
	RunLevel string `json:"run_level" textlog:"run_level"`
	// Privileges wanted by this scheduled task.
	Privileges StringList `json:"privileges,omitempty" textlog:"privileges,omitempty"`

	LastRun time.Time `json:"last_run,omitzero" textlog:"lastrun,omitempty"`
	NextRun time.Time `json:"next_run,omitzero" textlog:"nextrun,omitempty"`
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
	Guid        string    `json:"guid" textlog:"guid"`
	Path        string    `json:"path" textlog:"path"`
	Version     int       `json:"version" textlog:"version"`
	Created     time.Time `json:"created" textlog:"created"`
	LastRun     time.Time `json:"last_run" textlog:"last_run"`
	LastStopped time.Time `json:"last_stopped" textlog:"last_stopped"`
	Status      string    `json:"status" textlog:"status"`
	LastResult  string    `json:"last_result" textlog:"last_result"`
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

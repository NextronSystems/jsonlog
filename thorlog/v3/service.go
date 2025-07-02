package thorlog

import (
	"time"
)

type InitdService struct {
	LogObjectHeader

	File *File `json:"file" textlog:"file,expand"`
}

func (InitdService) reportable() {}

const typeInitdService = "init.d service"

func init() { AddLogObjectType(typeInitdService, &InitdService{}) }

func NewInitdService() *InitdService {
	return &InitdService{
		LogObjectHeader: LogObjectHeader{
			Type: typeInitdService,
		},
	}
}

type SystemdService struct {
	LogObjectHeader

	Command    string `json:"command" textlog:"command"`
	RunAsUser  string `json:"run_as_user" textlog:"run_as_user"`
	RunAsGroup string `json:"run_as_group" textlog:"run_as_group"`

	Unit  *File `json:"unit" textlog:"unit,expand"`
	Image *File `json:"image" textlog:"image,expand"`
}

func (SystemdService) reportable() {}

const typeSystemdService = "systemd service"

func init() { AddLogObjectType(typeSystemdService, &SystemdService{}) }

func NewSystemdService() *SystemdService {
	return &SystemdService{
		LogObjectHeader: LogObjectHeader{
			Type: typeSystemdService,
		},
	}
}

type WindowsService struct {
	LogObjectHeader

	Key            string    `json:"key" textlog:"key"`
	KeyName        string    `json:"key_name" textlog:"key_name"`
	ServiceName    string    `json:"service_name" textlog:"service_name"`
	Modified       time.Time `json:"modified" textlog:"modified"`
	StartType      string    `json:"start_type" textlog:"start_type"`
	ServiceType    string    `json:"service_type" textlog:"service_type"`
	User           string    `json:"user" textlog:"user"`
	Description    string    `json:"description" textlog:"description"`
	FailureCommand string    `json:"failure_command" textlog:"failure_command,omitempty"`
	Image          *File     `json:"image" textlog:"image,expand"`
}

func (WindowsService) reportable() {}

const typeWindowsService = "Windows service"

func init() { AddLogObjectType(typeWindowsService, &WindowsService{}) }

func NewWindowsService() *WindowsService {
	return &WindowsService{
		LogObjectHeader: LogObjectHeader{
			Type: typeWindowsService,
		},
	}
}

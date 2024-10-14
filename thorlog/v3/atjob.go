package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type AtJob struct {
	jsonlog.ObjectHeader

	File      string `json:"file" textlog:"job"`
	Command   string `json:"command" textlog:"command"`
	Start     string `json:"start" textlog:"start"`
	User      string `json:"user" textlog:"user"`
	RunLevel  string `json:"run_level" textlog:"runlevel"`
	LogonType string `json:"logon_type" textlog:"logontype"`
	Image     *File  `json:"image" textlog:"image,expand"`
}

const typeAtJob = "At Job"

func init() { AddLogObjectType(typeAtJob, &AtJob{}) }

func NewAtJob() *AtJob {
	return &AtJob{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeAtJob,
		},
	}
}

func (AtJob) reportable() {}

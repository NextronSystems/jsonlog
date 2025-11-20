package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type AtJob struct {
	jsonlog.ObjectHeader

	Command string `json:"command" textlog:"command"`
}

const typeAtJob = "at job"

func init() { AddLogObjectType(typeAtJob, &AtJob{}) }

func NewAtJob() *AtJob {
	return &AtJob{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeAtJob,
		},
	}
}

func (AtJob) reportable() {}

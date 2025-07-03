package thorlog

type CronJob struct {
	LogObjectHeader

	File     string `json:"file" textlog:"file"`
	User     string `json:"user" textlog:"user"`
	Schedule string `json:"schedule" textlog:"schedule"`
	Command  string `json:"command" textlog:"command"`
}

func (CronJob) reportable() {}

const typeCronJob = "cron job"

func init() { AddLogObjectType(typeCronJob, &CronJob{}) }

func NewCronjob() *CronJob {
	return &CronJob{
		LogObjectHeader: LogObjectHeader{
			Type: typeCronJob,
		},
	}
}

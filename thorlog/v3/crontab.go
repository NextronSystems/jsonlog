package thorlog

type CronJob struct {
	LogObjectHeader

	User     string `json:"user" textlog:"user"`
	Schedule string `json:"schedule" textlog:"schedule"`
	Command  string `json:"command" textlog:"command"`
}

func (CronJob) observed() {}

const typeCronJob = "cron job"

func init() { AddLogObjectType(typeCronJob, &CronJob{}) }

func NewCronjob() *CronJob {
	return &CronJob{
		LogObjectHeader: LogObjectHeader{
			Type: typeCronJob,
		},
	}
}

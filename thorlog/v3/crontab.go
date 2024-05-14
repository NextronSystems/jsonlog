package thorlog

type Cronjob struct {
	LogObjectHeader

	File     string `json:"file" textlog:"file"`
	User     string `json:"user" textlog:"user"`
	Schedule string `json:"schedule" textlog:"schedule"`
	Command  string `json:"command" textlog:"command"`
}

const typeCronjob = "cronjob"

func init() { AddLogObjectType(typeCronjob, &Cronjob{}) }

func NewCronjob() *Cronjob {
	return &Cronjob{
		LogObjectHeader: LogObjectHeader{
			Type: typeCronjob,
		},
	}
}

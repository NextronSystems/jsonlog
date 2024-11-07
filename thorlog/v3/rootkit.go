package thorlog

type Rootkit struct {
	LogObjectHeader
}

const typeRootkit = "rootkit"

func (Rootkit) reportable() {}

func init() { AddLogObjectType(typeRootkit, &Rootkit{}) }

func NewRootkit() *Rootkit {
	return &Rootkit{
		LogObjectHeader: LogObjectHeader{
			Type: typeRootkit,
		},
	}
}

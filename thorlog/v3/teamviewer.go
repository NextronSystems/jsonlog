package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type TeamViewerPassword struct {
	jsonlog.ObjectHeader
	Password string `json:"password" textlog:"password"`
	Name     string `json:"name" textlog:"name"`
}

func (TeamViewerPassword) observed() {}

const typeTeamViewerPassword = "TeamViewer password"

func init() { AddLogObjectType(typeTeamViewerPassword, &TeamViewerPassword{}) }

func NewTeamViewerPassword() *TeamViewerPassword {
	return &TeamViewerPassword{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeTeamViewerPassword,
		},
	}
}

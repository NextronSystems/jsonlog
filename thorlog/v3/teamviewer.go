package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type TeamViewerPassword struct {
	jsonlog.ObjectHeader
	Password     string `json:"password" textlog:"password"`
	Path         string `json:"path" textlog:"path"`
	RegistryPath string `json:"registry_path" textlog:"registry_path"`
	Name         string `json:"name" textlog:"name"`
}

func (TeamViewerPassword) reportable() {}

const typeTeamViewerPassword = "TeamViewer password"

func init() { AddLogObjectType(typeTeamViewerPassword, &TeamViewerPassword{}) }

func NewTeamViewerPassword() *TeamViewerPassword {
	return &TeamViewerPassword{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeTeamViewerPassword,
		},
	}
}

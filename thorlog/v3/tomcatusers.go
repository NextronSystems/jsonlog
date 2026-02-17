package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type TomcatUser struct {
	jsonlog.ObjectHeader
	User string `json:"user" textlog:"user"`
}

func (TomcatUser) observed() {}

const typeTomcatUser = "Tomcat user"

func init() { AddLogObjectType(typeTomcatUser, &TomcatUser{}) }

func NewTomcatUser(user string) *TomcatUser {
	return &TomcatUser{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeTomcatUser,
		},
		User: user,
	}
}

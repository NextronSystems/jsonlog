package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type TomcatUser struct {
	jsonlog.ObjectHeader
	User string `json:"user" textlog:"user"`
	File string `json:"file" textlog:"file"`
}

func (TomcatUser) reportable() {}

const typeTomcatUser = "Tomcat user"

func init() { AddLogObjectType(typeTomcatUser, &TomcatUser{}) }

func NewTomcatUser(user, file string) *TomcatUser {
	return &TomcatUser{
		ObjectHeader: jsonlog.ObjectHeader{
			Summary: "User " + user,
			Type:    typeTomcatUser,
		},
		User: user,
		File: file,
	}
}

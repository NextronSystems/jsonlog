package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type GroupsXmlUser struct {
	jsonlog.ObjectHeader
	User     string `json:"user" textlog:"user"`
	Password string `json:"password" textlog:"password"`
}

func (GroupsXmlUser) observed() {}

const typeGroupsXmlPassword = "groups.xml user"

func init() { AddLogObjectType(typeGroupsXmlPassword, &GroupsXmlUser{}) }

func NewGroupsXmlPassword(user, password string) *GroupsXmlUser {
	return &GroupsXmlUser{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeGroupsXmlPassword,
		},
		User:     user,
		Password: password,
	}
}

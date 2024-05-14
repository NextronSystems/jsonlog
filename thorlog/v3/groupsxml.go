package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type GroupsXmlUser struct {
	jsonlog.ObjectHeader
	File     string `json:"file" textlog:"file"`
	User     string `json:"user" textlog:"user"`
	Password string `json:"password" textlog:"password"`
}

const typeGroupsXmlPassword = "groups.xml user"

func init() { AddLogObjectType(typeGroupsXmlPassword, &GroupsXmlUser{}) }

func NewGroupsXmlPassword(file, user, password string) *GroupsXmlUser {
	return &GroupsXmlUser{
		ObjectHeader: jsonlog.ObjectHeader{
			Type:    typeGroupsXmlPassword,
			Summary: user,
		},
		File:     file,
		User:     user,
		Password: password,
	}
}

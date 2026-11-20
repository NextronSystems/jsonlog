package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type GroupsXmlUser struct {
	jsonlog.ObjectHeader
	User     string `json:"user" textlog:"user"`
	Password string `json:"password" textlog:"password"`
	// Changed is the time when this entry was last modified, taken from the
	// "changed" attribute of the corresponding element in the groups.xml file.
	Changed time.Time `json:"changed,omitzero" textlog:"changed,omitempty"`
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

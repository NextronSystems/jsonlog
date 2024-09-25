package thorlog

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/NextronSystems/jsonlog"
)

type Permissions interface {
	jsonlog.Object
	isPermission()
	String() string
}

var (
	_ Permissions = UnixPermissions{}
	_ Permissions = WindowsPermissions{}
)

func (UnixPermissions) isPermission()    {}
func (WindowsPermissions) isPermission() {}

type UnixPermissions struct {
	LogObjectHeader

	Owner string         `json:"owner" textlog:"owner"` // FIXME: Could explicitly include name / UID
	Group string         `json:"group" textlog:"group"` // FIXME: Could explicitly include name / GID
	Mask  PermissionMask `json:"permissions" textlog:"permissions"`
}

func (p UnixPermissions) String() string {
	return p.Mask.String()
}

type PermissionMask struct {
	User  RwxPermissions `json:"user"`
	Group RwxPermissions `json:"group"`
	World RwxPermissions `json:"world"`
}

func (p PermissionMask) String() string {
	return p.User.String() + p.Group.String() + p.World.String()
}

type RwxPermissions struct {
	Readable   bool `json:"readable"`
	Writable   bool `json:"writable"`
	Executable bool `json:"executable"`
}

func (r RwxPermissions) String() string {
	var s strings.Builder
	s.Grow(3)
	if r.Readable {
		s.WriteByte('r')
	} else {
		s.WriteByte('-')
	}
	if r.Writable {
		s.WriteByte('w')
	} else {
		s.WriteByte('-')
	}
	if r.Executable {
		s.WriteByte('x')
	} else {
		s.WriteByte('-')
	}
	return s.String()
}

const typeUnixPermissions = "unix permissions"

func init() { AddLogObjectType(typeUnixPermissions, &UnixPermissions{}) }

func NewUnixPermissions() *UnixPermissions {
	return &UnixPermissions{
		LogObjectHeader: jsonlog.ObjectHeader{
			Type: typeUnixPermissions,
		},
	}
}

type WindowsPermissions struct {
	LogObjectHeader

	Owner       string     `json:"owner" textlog:"owner"` // FIXME: Could include information like the original SID
	Permissions AclEntries `json:"permissions" textlog:"permissions"`
}

func (p WindowsPermissions) String() string {
	return p.Permissions.String()
}

type AclEntries []AclEntry

func (a AclEntries) String() string {
	var entryStrings = make([]string, len(a))
	for i, entry := range a {
		entryStrings[i] = entry.String()
	}
	sort.Strings(entryStrings)
	return strings.Join(entryStrings, " / ")
}

type AclEntry struct {
	Group  string    // FIXME: Could include information like the original SID
	Access AclAccess // FIXME: Could include the full original byte mask
}

func (a AclEntry) String() string {
	return a.Group + ":" + a.Access.String()
}

type AclAccess byte

const (
	FullPerm    AclAccess = 'F'
	ChangePerm  AclAccess = 'C'
	WritePerm   AclAccess = 'W'
	ReadPerm    AclAccess = 'R'
	SpecialPerm AclAccess = 'S'
)

func (a AclAccess) String() string {
	return string(a)
}

func (a AclAccess) MarshalJSON() ([]byte, error) {
	return []byte(`"` + a.String() + `"`), nil
}

func (a *AclAccess) UnmarshalJSON(data []byte) error {
	var accessString string
	if err := json.Unmarshal(data, &accessString); err != nil {
		return err
	}
	if len(accessString) != 1 {
		return fmt.Errorf("invalid ACL access string %s", accessString)
	}
	*a = AclAccess(accessString[0])
	return nil
}

const typeWindowsPermissions = "windows permissions"

func init() { AddLogObjectType(typeWindowsPermissions, &WindowsPermissions{}) }

func NewWindowsPermissions() *WindowsPermissions {
	return &WindowsPermissions{
		LogObjectHeader: jsonlog.ObjectHeader{
			Type: typeWindowsPermissions,
		},
	}
}

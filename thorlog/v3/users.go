package thorlog

import (
	"fmt"
	"time"
)

type LoggedInUser struct {
	LogObjectHeader

	User string `json:"user" textlog:"user"`

	Server       string `json:"server,omitempty" textlog:"server,omitempty"`
	Domain       string `json:"domain,omitempty" textlog:"domain,omitempty"`
	OtherDomains string `json:"other_domains,omitempty" textlog:"other_domains,omitempty"`
}

const typeLoggedInUser = "logged in user"

func init() { AddLogObjectType(typeLoggedInUser, &LoggedInUser{}) }

func NewLoggedInUser(user string) *LoggedInUser {
	return &LoggedInUser{
		LogObjectHeader: LogObjectHeader{
			Type:    typeLoggedInUser,
			Summary: user,
		},
		User: user,
	}
}

type ProfileFolder struct {
	LogObjectHeader

	User string `json:"user" textlog:"user"`

	Modified time.Time  `json:"modified" textlog:"modified,omitempty"`
	Created  *time.Time `json:"created,omitempty" textlog:"created,omitempty"`
}

const typeUserProfile = "user profile"

func init() { AddLogObjectType(typeUserProfile, &ProfileFolder{}) }

func NewProfileFolder(user string) *ProfileFolder {
	return &ProfileFolder{
		LogObjectHeader: LogObjectHeader{
			Type:    typeUserProfile,
			Summary: user,
		},
		User: user,
	}
}

type UnixUser struct {
	LogObjectHeader

	Name        string   `json:"name" textlog:"name"`
	Uid         string   `json:"uid" textlog:"userid"`
	Gid         string   `json:"gid" textlog:"groupid"`
	FullName    string   `json:"full_name" textlog:"full_name"`
	Home        string   `json:"home" textlog:"home"`
	Shell       string   `json:"shell" textlog:"shell"`
	Crontab     string   `json:"crontab" textlog:"-"`
	AccessFiles []string `json:"access_files" textlog:"-"`
}

const typeUnixUser = "unix user"

func init() { AddLogObjectType(typeUnixUser, &UnixUser{}) }

func NewUnixUser(name string) *UnixUser {
	return &UnixUser{
		LogObjectHeader: LogObjectHeader{
			Type:    typeUnixUser,
			Summary: name,
		},
		Name: name,
	}
}

type WindowsUser struct {
	LogObjectHeader

	User                 string       `json:"user" textlog:"user"`
	FullName             string       `json:"full_name" textlog:"full_name"`
	IsAdmin              bool         `json:"is_admin" textlog:"is_admin"`
	LastLogon            time.Time    `json:"last_logon" textlog:"last_logon"`
	BadPasswordCount     int          `json:"bad_password_count" textlog:"bad_password_count"`
	NumberOfLogons       int          `json:"num_logons" textlog:"num_logons"`
	PasswordAge          HourDuration `json:"pass_age" textlog:"pass_age"`
	PasswordNeverExpires bool         `json:"no_expire" textlog:"no_expire"`
	IsEnabled            bool         `json:"active" textlog:"active"`
	IsLocked             bool         `json:"locked" textlog:"locked"`
	Comment              string       `json:"comment" textlog:"comment"`
}

const typeWindowsUser = "windows user"

func init() { AddLogObjectType(typeWindowsUser, &WindowsUser{}) }

func NewWindowsUser(user string) *WindowsUser {
	return &WindowsUser{
		LogObjectHeader: LogObjectHeader{
			Type:    typeWindowsUser,
			Summary: user,
		},
		User: user,
	}
}

type HourDuration time.Duration

func (h HourDuration) String() string {
	return fmt.Sprintf("%.2f", time.Duration(h).Hours()/24)
}

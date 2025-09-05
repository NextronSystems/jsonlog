package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

// SRUMEntry holds information about a single entry of a System Resource Usage Monitor (SRUM)
// database. These databases are written by the SRUM service which collects and aggregates
// system resource usage data such as network activity, energy consumption, and application usage.
//
// Reference: https://www.forensafe.com/blogs/srudb.html
//
// A SRUMEntry represents a single entry in the "Application Resource Usage" table.
// Enriched with AppInfo, UserSID and UserName from the "SruDbIdMapTable" table.
//
// Columns in {D10CA2FE-6FCF-4F6D-848E-B2E99266FA89} (19):
// Id     Name                                 Type
// 2      TimeStamp                            DateTime
// 3      AppId                                Signed long
// 4      UserId                               Signed long
// 5      ForegroundCycleTime                  Long long
// 6      BackgroundCycleTime                  Long long
// 7      FaceTime                             Long long
// 10     ForegroundBytesRead                  Long long
// 11     ForegroundBytesWritten               Long long
// 15     BackgroundBytesRead                  Long long
// 16     BackgroundBytesWritten               Long long
//
// Columns in SruDbIdMapTable (3):
// Id     Name                                 Type
// 1      IdType                               Signed byte
// 2      IdIndex                              Signed long
// 256    IdBlob                               Long Binary

type SRUMEntry struct {
	jsonlog.ObjectHeader

	// TimeStamp is when the entry was recorded.
	TimeStamp time.Time `json:"timestamp" textlog:"timestamp"`

	// AppId is the numeric identifier of the application.
	AppId uint32 `json:"app_id" textlog:"app_id"`

	// AppInfo is the Appname/Apppath decoded from the IdBlob
	AppInfo string `json:"app_info" textlog:"app_info"`

	// UserId is the numeric identifier of the user.
	UserId uint32 `json:"user_id" textlog:"user_id"`

	// UserSID is the string SID parsed from the binary SID
	UserSID string `json:"user_info" textlog:"user_info"`

	// UserName is the Username looked up from the SID
	UserName string `json:"user_name" textlog:"user_name"`

	// FaceTime is the total foreground time in milliseconds.
	FaceTime uint64 `json:"face_time" textlog:"face_time"`

	// ForegroundBytesRead is the number of bytes read in the foreground.
	ForegroundBytesRead uint64 `json:"foreground_bytes_read" textlog:"foreground_bytes_read"`

	// ForegroundBytesWritten is the number of bytes written in the foreground.
	ForegroundBytesWritten uint64 `json:"foreground_bytes_written" textlog:"foreground_bytes_written"`

	// BackgroundBytesRead is the number of bytes read in the background.
	BackgroundBytesRead uint64 `json:"background_bytes_read" textlog:"background_bytes_read"`

	// BackgroundBytesWritten is the number of bytes written in the background.
	BackgroundBytesWritten uint64 `json:"background_bytes_written" textlog:"background_bytes_written"`
}

const typeSRUMEntry = "System Resource Usage Monitor"

func init() { AddLogObjectType(typeSRUMEntry, &SRUMEntry{}) }

func NewSRUMEntry() *SRUMEntry {
	return &SRUMEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeSRUMEntry,
		},
	}
}

func (SRUMEntry) reportable() {}

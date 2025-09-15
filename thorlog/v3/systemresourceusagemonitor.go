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
// 7      FaceTime                             Long long
// 10     ForegroundBytesRead                  Long long
// 11     ForegroundBytesWritten               Long long
// 12     ForegroundNumReadOperations          Long long
// 13     ForegroundNumWriteOperations         Long long
// 15     BackgroundBytesRead                  Long long
// 16     BackgroundBytesWritten               Long long
// 17     BackgroundNumReadOperations          Long long
// 18     BackgroundNumWriteOperations         Long long
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

	// AppInfo is the Appname/Apppath decoded from the IdBlob
	AppInfo string `json:"app_info" textlog:"app_info"`

	// UserSID is the string SID parsed from the binary SID
	UserSID string `json:"user_sid" textlog:"user_sid"`

	// UserName is the Username looked up from the SID, works only on Windows Builds.
	UserName string `json:"user_name" textlog:"user_name"`

	// FaceTime is the total foreground time in milliseconds.
	FaceTime uint64 `json:"face_time" textlog:"face_time"`

	// ForegroundBytesRead is the number of bytes read in the foreground.
	ForegroundBytesRead uint64 `json:"foreground_bytes_read" textlog:"foreground_bytes_read"`

	// ForegroundBytesWritten is the number of bytes written in the foreground.
	ForegroundBytesWritten uint64 `json:"foreground_bytes_written" textlog:"foreground_bytes_written"`

	// ForegroundNumReadOperations is the number of read operations in the foreground.
	ForegroundNumReadOperations uint64 `json:"foreground_num_read_operations" textlog:"foreground_num_read_operations"`

	// ForegroundNumWriteOperations is the number of write operations in the foreground.
	ForegroundNumWriteOperations uint64 `json:"foreground_num_write_operations" textlog:"foreground_num_write_operations"`

	// BackgroundBytesRead is the number of bytes read in the background.
	BackgroundBytesRead uint64 `json:"background_bytes_read" textlog:"background_bytes_read"`

	// BackgroundBytesWritten is the number of bytes written in the background.
	BackgroundBytesWritten uint64 `json:"background_bytes_written" textlog:"background_bytes_written"`

	// BackgroundNumReadOperations is the number of read operations in the background.
	BackgroundNumReadOperations uint64 `json:"background_num_read_operations" textlog:"background_num_read_operations"`

	// BackgroundNumWriteOperations is the number of write operations in the background.
	BackgroundNumWriteOperations uint64 `json:"background_num_write_operations" textlog:"background_num_write_operations"`
}

const typeSRUMEntry = "SRUM Resource Usage Entry"

func init() { AddLogObjectType(typeSRUMEntry, &SRUMEntry{}) }

func NewSRUMEntry() *SRUMEntry {
	return &SRUMEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeSRUMEntry,
		},
	}
}

func (SRUMEntry) reportable() {}

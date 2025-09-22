package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

// SRUMResourceUsageEntry holds information about a single entry of a System Resource Usage Monitor (SRUM)
// database. These databases are written by the SRUM service which collects and aggregates
// system resource usage data such as network activity, energy consumption, and application usage.
//
// Reference: https://www.forensafe.com/blogs/srudb.html
//
// A SRUMResourceUsageEntry represents a single entry in the "Application Resource Usage" table
// ({D10CA2FE-6FCF-4F6D-848E-B2E99266FA89}) enriched with AppInfo, UserSID and UserName
// from the "SruDbIdMapTable" table. Each entry represents a snapshot of resource usage
// for a specific application and user combination at a given time.
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
type SRUMResourceUsageEntry struct {
	jsonlog.ObjectHeader

	// TimeStamp is when the resource usage measurement was recorded by SRUM.
	// This represents the end time of the measurement period (typically hourly).
	TimeStamp time.Time `json:"timestamp" textlog:"timestamp"`

	// AppInfo contains the application path or executable name extracted from the
	// SruDbIdMapTable.IdBlob field. This identifies which application the resource
	// usage data belongs to (e.g., "C:\Windows\System32\notepad.exe").
	AppInfo string `json:"app_info" textlog:"app_info"`

	// UserSID is the Windows Security Identifier string parsed from the binary SID
	// stored in SruDbIdMapTable.IdBlob. This identifies which user account was
	// running the application (e.g., "S-1-5-21-...").
	UserSID string `json:"user_sid" textlog:"user_sid"`

	// UserName is the human-readable username resolved from the UserSID.
	// May be empty if the SID cannot be resolved to a username.
	UserName string `json:"user_name,omitempty" textlog:"user_name,omitempty"`

	// FaceTime is the total time in milliseconds that the application was visible
	// to the user (in the foreground) during the measurement period. This indicates
	// actual user interaction time with the application.
	FaceTime uint64 `json:"face_time" textlog:"face_time"`

	// ForegroundBytesRead is the total number of bytes read from disk/storage
	// while the application was in the foreground during the measurement period.
	ForegroundBytesRead uint64 `json:"foreground_bytes_read" textlog:"foreground_bytes_read"`

	// ForegroundBytesWritten is the total number of bytes written to disk/storage
	// while the application was in the foreground during the measurement period.
	ForegroundBytesWritten uint64 `json:"foreground_bytes_written" textlog:"foreground_bytes_written"`

	// ForegroundNumReadOperations is the count of discrete read I/O operations
	// performed while the application was in the foreground. This differs from
	// bytes read as it counts individual operations regardless of size.
	ForegroundNumReadOperations uint64 `json:"foreground_num_read_operations" textlog:"foreground_num_read_operations"`

	// ForegroundNumWriteOperations is the count of discrete write I/O operations
	// performed while the application was in the foreground. This differs from
	// bytes written as it counts individual operations regardless of size.
	ForegroundNumWriteOperations uint64 `json:"foreground_num_write_operations" textlog:"foreground_num_write_operations"`

	// BackgroundBytesRead is the total number of bytes read from disk/storage
	// while the application was running in the background during the measurement period.
	BackgroundBytesRead uint64 `json:"background_bytes_read" textlog:"background_bytes_read"`

	// BackgroundBytesWritten is the total number of bytes written to disk/storage
	// while the application was running in the background during the measurement period.
	BackgroundBytesWritten uint64 `json:"background_bytes_written" textlog:"background_bytes_written"`

	// BackgroundNumReadOperations is the count of discrete read I/O operations
	// performed while the application was running in the background. This differs
	// from bytes read as it counts individual operations regardless of size.
	BackgroundNumReadOperations uint64 `json:"background_num_read_operations" textlog:"background_num_read_operations"`

	// BackgroundNumWriteOperations is the count of discrete write I/O operations
	// performed while the application was running in the background. This differs
	// from bytes written as it counts individual operations regardless of size.
	BackgroundNumWriteOperations uint64 `json:"background_num_write_operations" textlog:"background_num_write_operations"`
}

const typeSRUMEntry = "SRUM Resource Usage Entry"

func init() { AddLogObjectType(typeSRUMEntry, &SRUMResourceUsageEntry{}) }

func NewSRUMResourceUsageEntry() *SRUMResourceUsageEntry {
	return &SRUMResourceUsageEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeSRUMEntry,
		},
	}
}

func (SRUMResourceUsageEntry) reportable() {}

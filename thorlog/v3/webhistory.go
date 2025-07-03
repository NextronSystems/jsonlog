package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

// WebDownload is a log object that represents a web download.
//
// The download is not guaranteed to be complete or successful.
type WebDownload struct {
	jsonlog.ObjectHeader

	// URL is the URL of the downloaded file.
	URL string `json:"url" textlog:"url"`
	// Time is the time when the download was started.
	Time time.Time `json:"time" textlog:"time"`

	// File contains the information about the downloaded file.
	File *File `json:"file" textlog:"file,expand"`
}

func (WebDownload) reportable() {}

const typeWebDownload = "web download"

func init() { AddLogObjectType(typeWebDownload, &WebDownload{}) }

func NewWebDownload() *WebDownload {
	return &WebDownload{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeWebDownload,
		},
	}
}

// WebPageVisit is a log object that represents a web page visit.
//
// The visit may also have been triggered indirectly (e.g. a JavaScript file that was loaded).
type WebPageVisit struct {
	jsonlog.ObjectHeader

	URL  string    `json:"url" textlog:"url"`
	Time time.Time `json:"time" textlog:"time"`
	// Title is the title of the visited page.
	Title string `json:"title" textlog:"title"`
}

func (WebPageVisit) reportable() {}

const typeWebVisit = "web page visit"

func init() { AddLogObjectType(typeWebVisit, &WebPageVisit{}) }

func NewWebVisit() *WebPageVisit {
	return &WebPageVisit{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeWebVisit,
		},
	}
}

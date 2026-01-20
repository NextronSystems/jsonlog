package thorlog

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/thorlog/common"
)

func TestContext_MarshalTextLog(t *testing.T) {
	tests := []struct {
		name string
		c    *Context
		want string
	}{
		{
			name: "empty context",
			c:    &Context{},
			want: "",
		},
		{
			name: "context with unique related object",
			c: &Context{
				{
					Object: NewFile("path/to/file"),
					Relations: []Relation{{
						Name:   "file",
						Unique: true,
					}},
				},
			},
			want: "FILE: path/to/file",
		},
		{
			name: "context with related object group",
			c: &Context{
				{
					Object: NewFile("path/to/file"),
					Relations: []Relation{{
						Name:   "file",
						Unique: false,
					}},
				},
				{
					Object: NewFile("path/to/otherfile"),
					Relations: []Relation{{
						Name:   "file",
						Unique: false,
					}},
				},
			},
			want: "FILE_1: path/to/file FILE_2: path/to/otherfile",
		},
		{
			name: "context with different related objects",
			c: &Context{
				{
					Object: NewFile("path/to/file"),
					Relations: []Relation{{
						Name:   "file",
						Unique: false,
					}},
				},
				{
					Object: NewFile("path/to/otherfile"),
					Relations: []Relation{{
						Name:   "archive",
						Unique: true,
					}},
				},
			},
			want: "FILE_1: path/to/file ARCHIVE_FILE: path/to/otherfile",
		},
		{
			name: "context with object related in two ways",
			c: &Context{
				{
					Object: NewFile("path/to/file"),
					Relations: []Relation{{
						Name:   "parent",
						Type:   "derived from",
						Unique: true,
					}, {
						Name:   "origin",
						Type:   "derived from",
						Unique: true,
					}},
				},
			},
			want: "PARENT_FILE: path/to/file",
		},
	}
	var formatter jsonlog.TextlogFormatter
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := concatEntry(formatter.Format(tt.c)); got != tt.want {
				t.Errorf("Context.MarshalTextLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func concatEntry(entry jsonlog.TextlogEntry) string {
	var builder strings.Builder
	for i, e := range entry {
		if i > 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(e.Key)
		builder.WriteString(": ")
		builder.WriteString(e.Value)
	}
	return builder.String()
}

func TestFinding_UnmarshalJSON(t *testing.T) {
	for i, finding := range []*Finding{
		{
			ObjectHeader: LogObjectHeader{Type: typeFinding},
			Meta: LogEventMetadata{
				Lvl:    common.Alert,
				Mod:    "Test",
				ScanID: "abdc",
				GenID:  "abdas",
				Source: "aserarsd",
			},
			Text:    "This is a test finding",
			Subject: NewFile("path/to/file"),
			EventContext: Context{
				{
					Object: NewAtJob(),
					Relations: []Relation{{
						Type: "related to",
					}},
				},
			},
			Reasons: []Reason{
				NewReason("Reason 1", Signature{Score: 70}, nil),
			},
			Score:      70,
			LogVersion: common.JsonV3,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			jsonform, err := json.Marshal(finding)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(jsonform))
			var newFinding Finding
			if err := json.Unmarshal(jsonform, &newFinding); err != nil {
				t.Fatal(err)
			}
			newFinding.LogVersion = common.Version(newFinding.LogVersion.Major())
			if !reflect.DeepEqual(finding, &newFinding) {
				t.Errorf("UnmarshalJSON() = %+v, want %+v", newFinding, *finding)
			}
		})
	}
}

func TestFinding_UnmarshalIssue(t *testing.T) {
	finding := `{"type":"THOR finding","meta":{"time":"2025-07-01T12:05:12.993789131+02:00","level":"Info","module":"ProcessCheck","scan_id":"S-pSxgCmyvvfs","event_id":"","hostname":"dummy"},"message":"process found","subject":{"type":"process","pid":502168,"name":"chromium","command":"/usr/lib/chromium/chromium","owner":"owner","image":{"type":"file","path":"/usr/lib/chromium/chromium","exists":"yes","extension":"","magic_header":"ELF","hashes":{"md5":"fc04ee20f064adc18e370c22512e268e","sha1":"2c8b7d05d25e04db9c169ce85e8e8f84321ef0c8","sha256":"0cf1727aa8dc3995d5aa103001f656b8ee8a1b3ffbc6d8664c5ad95cf225771f"},"first_bytes":{"hex":"7f454c4602010100000000000000000003003e00","ascii":"ELF\u003e"},"file_times":{"modified":"2025-06-25T19:45:43+02:00","accessed":"2025-07-01T08:46:56.750309598+02:00","changed":"2025-06-26T08:39:59.980605063+02:00"},"size":252546120,"permissions":{"type":"Unix permissions","owner":"root","group":"root","mask":{"user":{"readable":true,"writable":true,"executable":true},"group":{"readable":true,"writable":false,"executable":true},"world":{"readable":true,"writable":false,"executable":true}}}},"parent_info":{"pid":9011,"exe":"/usr/lib/chromium/chromium","command":"/usr/lib/chromium/chromium"},"tree":["/usr/lib/chromium/chromium","/usr/lib/chromium/chromium"],"created":"2025-07-01T12:00:05+02:00","session":"","listen_ports":null,"connections":[]},"score":0,"reasons":null,"reason_count":0,"context":null,"issues":[{"affected":"/subject/sections","category":"truncated","description":"Removed some sections from process memory (originally 638)"}],"log_version":"v3.0.0"}`
	var findingObj Finding
	if err := json.Unmarshal([]byte(finding), &findingObj); err != nil {
		t.Fatalf("Failed to unmarshal finding: %v", err)
	}
}

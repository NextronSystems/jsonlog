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
					Object:   NewFile("path/to/file"),
					Relation: "file",
					Unique:   true,
				},
			},
			want: "FILE: path/to/file",
		},
		{
			name: "context with related object group",
			c: &Context{
				{
					Object:   NewFile("path/to/file"),
					Relation: "file",
					Unique:   false,
				},
				{
					Object:   NewFile("path/to/otherfile"),
					Relation: "file",
					Unique:   false,
				},
			},
			want: "FILE_1: path/to/file FILE_2: path/to/otherfile",
		},
		{
			name: "context with different related objects",
			c: &Context{
				{
					Object:   NewFile("path/to/file"),
					Relation: "file",
					Unique:   false,
				},
				{
					Object:   NewFile("path/to/otherfile"),
					Relation: "archive",
					Unique:   true,
				},
			},
			want: "FILE_1: path/to/file ARCHIVE_FILE: path/to/otherfile",
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

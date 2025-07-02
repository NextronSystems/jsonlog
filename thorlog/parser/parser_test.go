package parser

import (
	"testing"
	"time"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/thorlog/common"
	thorlogv1 "github.com/NextronSystems/jsonlog/thorlog/v1"
	thorlogv2 "github.com/NextronSystems/jsonlog/thorlog/v2"
	"github.com/NextronSystems/jsonlog/thorlog/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestParseEvent(t *testing.T) {
	for _, testcase := range []struct {
		name     string
		rawEvent string
		expected thorlog.Event
	}{
		{
			"JsonV1Logline",
			`{"time":"2024-09-24T12:35:41Z","hostname":"host","level":"Info","module":"Startup","message":"Sigma Database: r2024-07-17-48-g5c4f599e3","scanid":"S-NtvqbRLWbf8","log_version":"v1.0.0"}`,
			&thorlogv1.Event{
				LogEventMetadata: thorlogv1.Metadata{
					Time:   mustTime("2024-09-24T12:35:41Z"),
					Lvl:    common.Info,
					Mod:    "Startup",
					ScanID: "S-NtvqbRLWbf8",
					GenID:  "",
					Source: "host",
				},
				Data: thorlogv1.Fields{
					{Key: "message", Value: "Sigma Database: r2024-07-17-48-g5c4f599e3"},
				},
			},
		},
		{
			"JsonV2",
			`{"time":"2024-09-24T12:37:04Z","hostname":"host","level":"Info","module":"Hosts","message":"Starting module","scanid":"S-kgKxYyJFQd0","log_version":"v2.0.0"}`,
			&thorlogv2.Event{
				LogEventMetadata: thorlogv2.Metadata{
					Time:   mustTime("2024-09-24T12:37:04Z"),
					Lvl:    common.Info,
					Mod:    "Hosts",
					ScanID: "S-kgKxYyJFQd0",
					GenID:  "",
					Source: "host",
				},
				Data: thorlogv2.Fields{
					{Key: "message", Value: "Starting module"},
				},
				EventVersion: "v2.0.0",
			},
		},
		{
			"JsonV2WithArray",
			`{"time":"2024-09-24T12:37:04Z","hostname":"host","level":"Info","module":"Hosts","message":"something","tags":["abc","def","ghi"],"scanid":"S-kgKxYyJFQd0","log_version":"v2.0.0"}`,
			&thorlogv2.Event{
				LogEventMetadata: thorlogv2.Metadata{
					Time:   mustTime("2024-09-24T12:37:04Z"),
					Lvl:    common.Info,
					Mod:    "Hosts",
					ScanID: "S-kgKxYyJFQd0",
					GenID:  "",
					Source: "host",
				},
				Data: thorlogv2.Fields{
					{Key: "message", Value: "something"},
					{Key: "tags", Value: []any{"abc", "def", "ghi"}},
				},
				EventVersion: "v2.0.0",
			},
		},
		{
			"JsonV2LoglineWithNested",
			`{"time":"2024-09-24T12:37:04Z","hostname":"host","level":"Info","module":"Hosts","message":"something","reasons":[{"reason":"r1","rulename":"rn1"},{"reason":"r2","rulename":"rn2"}],"scanid":"S-kgKxYyJFQd0","log_version":"v2.0.0"}`,
			&thorlogv2.Event{
				LogEventMetadata: thorlogv2.Metadata{
					Time:   mustTime("2024-09-24T12:37:04Z"),
					Lvl:    common.Info,
					Mod:    "Hosts",
					ScanID: "S-kgKxYyJFQd0",
					GenID:  "",
					Source: "host",
				},
				Data: thorlogv2.Fields{
					{Key: "message", Value: "something"},
					{Key: "reasons", Value: []any{
						thorlogv2.Fields{
							{Key: "reason", Value: "r1"},
							{Key: "rulename", Value: "rn1"},
						},
						thorlogv2.Fields{
							{Key: "reason", Value: "r2"},
							{Key: "rulename", Value: "rn2"},
						},
					}},
				},
				EventVersion: "v2.0.0",
			},
		},
		{
			"JsonV3Message",
			`{"message":"Starting module","type":"THOR message","meta":{"time":"2024-09-24T14:18:46.190394329+02:00","level":"Info","module":"Hosts","scan_id":"S-UBNfBD4xE8s","event_id":"","hostname":"host"},"fields":{},"log_version":"v3.0.0"}`,
			&thorlog.Message{
				ObjectHeader: jsonlog.ObjectHeader{
					Type: "THOR message",
				},
				Meta: thorlog.LogEventMetadata{
					Time:   mustTime("2024-09-24T14:18:46.190394329+02:00"),
					Lvl:    common.Info,
					Mod:    "Hosts",
					ScanID: "S-UBNfBD4xE8s",
					GenID:  "",
					Source: "host",
				},
				Text:       "Starting module",
				LogVersion: "v3.0.0",
			},
		},
		{
			"JsonV3Finding",
			`{"type":"THOR finding","meta":{"time":"2024-09-24T14:18:46.190394329+02:00","level":"Alert","module":"Test","scan_id":"abdc","event_id":"abdas","hostname":"aserarsd"},"message":"This is a test finding","subject":{"type":"file","path":"path/to/file"},"score":70,"reasons":[{"type":"reason","summary":"Reason 1","signature":{"score":70,"ref":null,"origin":"internal","kind":""},"matched":null}],"reason_count":0,"context":[{"object":{"type":"at job"},"relation":"","unique":false}],"log_version":"v3"}`,
			&thorlog.Finding{
				ObjectHeader: jsonlog.ObjectHeader{
					Type: "THOR finding",
				},
				Meta: thorlog.LogEventMetadata{
					Time:   mustTime("2024-09-24T14:18:46.190394329+02:00"),
					Lvl:    common.Alert,
					Mod:    "Test",
					ScanID: "abdc",
					GenID:  "abdas",
					Source: "aserarsd",
				},
				Text: "This is a test finding",
				Subject: &thorlog.File{
					ObjectHeader: jsonlog.ObjectHeader{
						Type: "file",
					},
					Path: "path/to/file",
				},
				Score: 70,
				Reasons: []thorlog.Reason{
					{
						ObjectHeader: jsonlog.ObjectHeader{
							Type: "reason",
						},
						Summary: "Reason 1",
						Signature: thorlog.Signature{
							Score: 70,
							Type:  thorlog.Internal,
						},
					},
				},
				ReasonCount: 0,
				EventContext: thorlog.Context{
					{
						Object: &thorlog.AtJob{
							ObjectHeader: jsonlog.ObjectHeader{
								Type: "at job",
							},
						},
					},
				},
				LogVersion: "v3.0.0",
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			event, err := ParseEvent([]byte(testcase.rawEvent))
			require.NoError(t, err)
			assert.Equal(t, testcase.expected, event)
		})
	}
}

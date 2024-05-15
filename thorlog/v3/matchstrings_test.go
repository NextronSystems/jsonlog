package thorlog

import (
	"testing"

	"github.com/NextronSystems/jsonlog"
)

func TestMatchString_String(t *testing.T) {
	var testObject = new(struct {
		jsonlog.ObjectHeader
		MyField string `json:"my_field" textlog:"my_field"`
	})
	testObject.MyField = "bar"
	var ref = jsonlog.NewReference(testObject, &testObject.MyField)
	for _, tt := range []struct {
		input    MatchString
		expected string
	}{
		{
			input: MatchString{
				Match: MatchData{
					Data: []byte("foo"),
				},
			},
			expected: "foo",
		},
		{
			input: MatchString{
				Match: MatchData{
					Data: []byte("foo"),
				},
				Context: &MatchData{
					Data: []byte("a wild foo appears"),
				},
			},
			expected: `foo in "a wild foo appears"`,
		},
		{
			input: MatchString{
				Match: MatchData{
					Data: []byte("foo"),
				},
				Context: &MatchData{
					Data: []byte("a wild foo appears"),
				},
			},
			expected: `foo in "a wild foo appears"`,
		},
		{
			input: MatchString{
				Match: MatchData{
					Data:    []byte("foo"),
					FullHex: true,
				},
			},
			expected: `666f6f`,
		},
		{
			input: MatchString{
				Match: MatchData{
					Data: []byte("foo"),
				},
				Offset: asOptional(0x10),
			},
			expected: `foo at 0x10`,
		},
		{
			input: MatchString{
				Match: MatchData{
					Data: []byte("foo"),
				},
				Offset:     asOptional(0x10),
				HideOffset: true,
			},
			expected: `foo`,
		},
		{
			input: MatchString{
				Match: MatchData{
					Data: []byte("foo"),
				},
				Offset: asOptional(0x10),
				Field:  ref,
			},
			expected: `foo at 0x10 in MY_FIELD`,
		},
		{
			input: MatchString{
				Match: MatchData{
					Data: []byte("bar"),
				},
				Offset: asOptional(0x10),
				Field:  ref,
			},
			expected: `bar at 0x10 in MY_FIELD`,
		},
		{
			input: MatchString{
				Match: MatchData{
					Data: []byte("bar"),
				},
				Offset: asOptional(0),
				Field:  ref,
			},
			expected: `bar in MY_FIELD`,
		},
		{
			input: MatchString{
				Match: MatchData{
					Data: []byte("fo\x00o"),
				},
			},
			expected: `"fo\x00o"`,
		},
	} {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.input.String(); got != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func asOptional(i uint64) *uint64 {
	return &i
}

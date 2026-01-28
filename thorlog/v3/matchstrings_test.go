package thorlog

import (
	"testing"

	"github.com/NextronSystems/jsonlog"
)

func TestEncodeString(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected StringWithEncoding
	}{
		{
			input: "hello",
			expected: StringWithEncoding{
				EncodedData: "hello",
				Encoding:    Plain,
			},
		},
		{
			input: string([]byte{0xff, 0xfe, 0xfd}),
			expected: StringWithEncoding{
				EncodedData: "fffefd",
				Encoding:    Hex,
			},
		},
	} {
		t.Run(tt.input, func(t *testing.T) {
			if got := EncodeString(tt.input); got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestStringWithEncoding_Data(t *testing.T) {
	for _, tt := range []struct {
		input    StringWithEncoding
		expected []byte
	}{
		{
			input: StringWithEncoding{
				EncodedData: "hello",
				Encoding:    Plain,
			},
			expected: []byte("hello"),
		},
		{
			input: StringWithEncoding{
				EncodedData: "68656c6c6f",
				Encoding:    Hex,
			},
			expected: []byte("hello"),
		},
	} {
		t.Run(tt.input.EncodedData, func(t *testing.T) {
			if got := tt.input.Data(); string(got) != string(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
func TestStringWithEncoding_String(t *testing.T) {
	for _, tt := range []struct {
		input    StringWithEncoding
		expected string
	}{
		{
			input: StringWithEncoding{
				EncodedData: "hello",
				Encoding:    Plain,
			},
			expected: "hello",
		},

		{
			input: StringWithEncoding{
				EncodedData: "666f6f",
				Encoding:    Hex,
			},
			expected: `foo`,
		},
		{
			input: StringWithEncoding{
				EncodedData: "010203",
				Encoding:    Hex,
			},
			expected: `"\x01\x02\x03"`,
		},
		{
			input: StringWithEncoding{
				EncodedData: "fo\x00o",
				Encoding:    Plain,
			},
			expected: `"fo\x00o"`,
		},
		{
			input: StringWithEncoding{
				EncodedData: `"quoted" data`,
				Encoding:    Plain,
			},
			expected: `"\"quoted\" data"`,
		},
		{
			input: StringWithEncoding{
				EncodedData: `8081fe`, // invalid UTF-8 byte sequence
				Encoding:    Hex,
			},
			expected: `"\x80\x81\xfe"`,
		},
		{
			input: StringWithEncoding{
				EncodedData: "a\x00b\x00c\x00", // UTF-16LE encoded "abc"
				Encoding:    Plain,
			},
			expected: `abc`,
		},
	} {
		t.Run(tt.input.EncodedData, func(t *testing.T) {
			if got := tt.input.String(); got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

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
				Match: EncodeString("foo"),
			},
			expected: "foo",
		},
		{
			input: MatchString{
				Match: EncodeString("foo"),
				Context: &StringWithEncoding{
					EncodedData: "a wild foo appears",
					Encoding:    Plain,
				},
			},
			expected: `foo in "a wild foo appears"`,
		},
		{
			input: MatchString{
				Match:  EncodeString("foo"),
				Offset: asOptional(0x10),
			},
			expected: `foo at 0x10`,
		},
		{
			input: MatchString{
				Match:      EncodeString("foo"),
				Offset:     asOptional(0x10),
				HideOffset: true,
			},
			expected: `foo`,
		},
		{
			input: MatchString{
				Match:  EncodeString("foo"),
				Offset: asOptional(0x10),
				Field:  ref,
			},
			expected: `foo at 0x10 in MY_FIELD`,
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

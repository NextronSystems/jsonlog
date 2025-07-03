package thorlog

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/NextronSystems/jsonlog"
)

// SparseData is a log object that represents a sparse data structure.
// SparseData represents a selection of data from a large data block (e.g.: a file's content)
// that is not fully contained in the log.
//
// Not all parts of the full data structure are necessarily contained in the sparse data,
// typically based on how much data is relevant for the analysis.
type SparseData struct {
	jsonlog.ObjectHeader
	// Elements is a list of sparse data elements that contain the actual data.
	// Each element has an offset within the block and the data that is present at that offset.
	// Elements are ordered by their offset, and are guaranteed to be non-overlapping.
	Elements []SparseDataElement `json:"elements" jsonschema:"nullable"`
	// Length is the length of the block where the sparse elements reside in.
	// In other words, all Elements are within an address range of [0, Length).
	Length int64 `json:"length"`
}

const truncateSequence = "[...]"

func (s SparseData) String() string {
	if s.Length == 0 {
		return ""
	}
	if len(s.Elements) == 0 {
		return truncateSequence
	}
	var w strings.Builder
	if s.Elements[0].Offset > 0 {
		_, _ = w.WriteString(truncateSequence)
	}
	for _, element := range s.Elements {
		_, _ = nonAsciiEscaper.WriteString(&w, string(element.Data))
		if element.Offset+uint64(len(element.Data)) < uint64(s.Length) {
			_, _ = w.WriteString(truncateSequence)
		}
	}
	return w.String()
}

type SparseDataElement struct {
	Offset uint64               `json:"offset"`
	Data   InvalidUnicodeString `json:"data"`
}

type InvalidUnicodeString string

func (s InvalidUnicodeString) MarshalJSON() ([]byte, error) {
	matchingString := escaper.Replace(string(s))
	var replacedString bytes.Buffer
	for _, char := range []byte(matchingString) {
		if char < 0x20 || char > 0x7E { // non ASCII
			replacedString.WriteString("\\u00")
			replacedString.WriteString(hex.EncodeToString([]byte{char}))
		} else {
			replacedString.WriteByte(char)
		}
	}
	matchingString = replacedString.String()
	matchingString = fmt.Sprintf("\"%s\"", matchingString)
	return []byte(matchingString), nil
}

var escaper = strings.NewReplacer("\\", "\\\\", "\"", "\\\"")

const typeSparseData = "sparse data"

func init() { AddLogObjectType(typeSparseData, &SparseData{}) }

func NewSparseData() *SparseData {
	return &SparseData{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeSparseData,
		},
	}
}

var nonAsciiEscaper *strings.Replacer

func init() {
	var escapes = []string{`\`, `\\`}
	for nonAsciiByte := 0; nonAsciiByte < 0x20; nonAsciiByte++ {
		escapes = append(escapes, string([]byte{byte(nonAsciiByte)}), fmt.Sprintf("\\x%02x", nonAsciiByte))
	}
	for nonAsciiByte := 0x7F; nonAsciiByte <= 0xFF; nonAsciiByte++ {
		escapes = append(escapes, string([]byte{byte(nonAsciiByte)}), fmt.Sprintf("\\x%02x", nonAsciiByte))
	}
	nonAsciiEscaper = strings.NewReplacer(escapes...)
}

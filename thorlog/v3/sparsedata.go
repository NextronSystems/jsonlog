package thorlog

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/NextronSystems/jsonlog"
)

type SparseData struct {
	jsonlog.ObjectHeader
	Elements []SparseDataElement `json:"elements"`
	Length   int64               `json:"length"`

	StringVersion string `json:"-" textlog:",explicit"`
}

func (s SparseData) String() string {
	return s.StringVersion
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

const typeSparseData = "sparsedata"

func init() { AddLogObjectType(typeSparseData, &SparseData{}) }

func NewSparseData() *SparseData {
	return &SparseData{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeSparseData,
		},
	}
}

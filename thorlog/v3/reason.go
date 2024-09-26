package thorlog

import (
	"encoding/json"
	"fmt"

	"github.com/NextronSystems/jsonlog"
)

type Reason struct {
	jsonlog.ObjectHeader

	Summary string `json:"-" textlog:"reason"`

	Signature     `json:"signature" textlog:",inline"`
	StringMatches MatchStrings `json:"matched" textlog:"matched"`
}

func (r *Reason) UnmarshalJSON(data []byte) error {
	type plainReason Reason
	if err := json.Unmarshal(data, (*plainReason)(r)); err != nil {
		return err
	}
	r.Summary = r.ObjectHeader.Summary
	return nil
}

const typeReason = "reason"

func init() {
	AddLogObjectType(typeReason, &Reason{})
}

type Signature struct {
	Ref             StringList `json:"ref" textlog:"ref"`
	Type            Sigtype    `json:"origin" textlog:"sigtype"`
	Class           Sigclass   `json:"kind" textlog:"sigclass"`
	Score           int64      `json:"score" textlog:"subscore"`
	Date            string     `json:"ruledate,omitempty" textlog:"ruledate,omitempty"`
	Tags            StringList `json:"tags,omitempty" textlog:"tags,omitempty"`
	Rulename        string     `json:"rulename,omitempty" textlog:"rulename,omitempty"`
	LongDescription string     `json:"description,omitempty" textlog:"description,omitempty"`
	Author          string     `json:"author,omitempty" textlog:"author,omitempty"`
	RuleId          string     `json:"id,omitempty" textlog:"id"`
	FalsePositives  StringList `json:"falsepositives,omitempty" textlog:"falsepositives,omitempty"`
}

type Sigclass string

const (
	ClassFilenameIOC       Sigclass = "Filename IOC"
	ClassNamedPipeIOC      Sigclass = "Named Pipe IOC"
	ClassYaraRule          Sigclass = "YARA Rule"
	ClassSigmaRule         Sigclass = "Sigma Rule"
	ClassStixIOC           Sigclass = "STIX IOC"
	ClassInternalHeuristic Sigclass = "Internal Heuristic"
	ClassHashIOC           Sigclass = "Hash IOC"
	ClassKeywordIOC        Sigclass = "Keyword IOC"
	ClassC2IOC             Sigclass = "Domain IOC"
	ClassHandleIOC         Sigclass = "Handle IOC"
)

type Sigtype int

const (
	Internal Sigtype = iota
	Custom
	External
)

func (s Sigtype) String() string {
	switch s {
	case Custom:
		return "custom"
	case Internal:
		return "internal"
	case External:
		return "external"
	default:
		return "unknown"
	}
}

func (s Sigtype) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Sigtype) UnmarshalJSON(data []byte) error {
	var stringType string
	if err := json.Unmarshal(data, &stringType); err != nil {
		return err
	}
	switch stringType {
	case "custom":
		*s = Custom
	case "internal":
		*s = Internal
	case "external":
		*s = External
	default:
		return fmt.Errorf("unknown sigtype %s", stringType)
	}
	return nil
}

func NewReason(desc string, signature Signature, matches MatchStrings) Reason {
	return Reason{
		ObjectHeader: jsonlog.ObjectHeader{
			Type:    typeReason,
			Summary: desc,
		},
		Summary:       desc,
		Signature:     signature,
		StringMatches: matches,
	}
}

package thorlog

import (
	"encoding/json"
	"fmt"

	"github.com/NextronSystems/jsonlog"
)

// Reason describes a match of a single Signature on a ObservedObject.
type Reason struct {
	jsonlog.ObjectHeader

	Summary string `json:"summary" textlog:"reason"`

	// Signature contains details about the signature that matched.
	Signature `json:"signature" textlog:",inline"`
	// StringMatches contains the matches that explain why this signature matched.
	StringMatches MatchStrings `json:"matched" textlog:"matched" jsonschema:"nullable"`
}

func (r *Reason) UnmarshalJSON(data []byte) error {
	type plainReason Reason
	if err := json.Unmarshal(data, (*plainReason)(r)); err != nil {
		return err
	}
	return nil
}

const typeReason = "reason"

func init() {
	AddLogObjectType(typeReason, &Reason{})
}

// Signature describes metadata about a signature that THOR uses to detect
// suspicious objects.
type Signature struct {
	// Score is a metric that combines severity and certainty for this signature.
	//
	// It is related to the Assessment.Score, which is derived from the scores of all
	// signatures that matched; however, signature scores are not limited to the
	// 0 to 100 interval of assessment scores, but may also be negative to indicate
	// a likely false positive (which results in a score reduction on any related
	// assessment).
	Score int64 `json:"score" textlog:"subscore"`
	// Ref contains references (usually as links) for further information about
	// the threat that is detected by this signature.
	Ref StringList `json:"reference" textlog:"ref" jsonschema:"nullable"`
	// Type indicates whether a signature was part of THOR's built in signature set
	// or whether it was a custom signature provided by the user.
	Type Sigtype `json:"origin" textlog:"sigtype"`
	// Class is the sort of signature that this is (YARA Rule, Filename IOC, ...)
	Class Sigclass `json:"kind" textlog:"sigclass"`
	// Date is the date on which the signature was last modified.
	Date string `json:"date,omitempty" textlog:"ruledate,omitempty"`
	// Tags are short strings that help with grouping signatures.
	//
	// E.g. APT related signatures may be tagged "APT", or malware related signatures may be tagged "MAL".
	Tags StringList `json:"tags,omitempty" textlog:"tags,omitempty" jsonschema:"nullable"`
	// Rulename is the name of the signature (e.g. a YARA rule name).
	Rulename string `json:"rule_name,omitempty" textlog:"rulename,omitempty"`
	// LongDescription contains the description that the signature has about itself
	// (e.g. "detects a webshell related to ...")
	LongDescription string `json:"description,omitempty" textlog:"description,omitempty"`
	// Author is the name of the person who wrote the signature.
	Author string `json:"author,omitempty" textlog:"author,omitempty"`
	// RuleId is a unique ID that identifies this signature.
	//
	// Not all classes of signatures may provide this field.
	RuleId string `json:"id,omitempty" textlog:"id"`
	// FalsePositives describes cases where this signature is known to produce matches
	// even on benign data.
	FalsePositives StringList `json:"false_positives,omitempty" textlog:"falsepositives,omitempty" jsonschema:"nullable"`
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

func (s Sigtype) JSONSchemaAlias() any { return "" }

func NewReason(desc string, signature Signature, matches MatchStrings) Reason {
	return Reason{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeReason,
		},
		Summary:       desc,
		Signature:     signature,
		StringMatches: matches,
	}
}

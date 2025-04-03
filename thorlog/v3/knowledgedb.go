package thorlog

import (
	"time"
)

type KnowledgeDBEntry struct {
	LogObjectHeader

	Entry      string        `json:"entry" textlog:"entry"`
	Created    Time          `json:"created" textlog:"created"`
	Started    Time          `json:"started" textlog:"started"`
	Duration   time.Duration `json:"duration" textlog:"duration"`
	PrimaryKey int64         `json:"primary_key" textlog:"primary_key"`
}

func (KnowledgeDBEntry) reportable() {}

const typeKnowledgeDBEntry = "KnowledgeDB entry"

func init() { AddLogObjectType(typeKnowledgeDBEntry, &KnowledgeDBEntry{}) }

func NewKnowledgeDBEntry() *KnowledgeDBEntry {
	return &KnowledgeDBEntry{
		LogObjectHeader: LogObjectHeader{
			Type: typeKnowledgeDBEntry,
		},
	}
}

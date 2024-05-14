package thorlog

import (
	"fmt"
)

type ThunderstormSampleContext struct {
	LogObjectHeader
	SampleId int64 `json:"sample_id" textlog:"sample_id"`
}

const typeThunderstormSampleContext = "Thunderstorm context"

func init() { AddLogObjectType(typeThunderstormSampleContext, &ThunderstormSampleContext{}) }

func NewThunderstormSampleContext(sampleId int64) *ThunderstormSampleContext {
	return &ThunderstormSampleContext{
		LogObjectHeader: LogObjectHeader{
			Summary: fmt.Sprintf("event happened during scan of sample %d", sampleId),
			Type:    typeThunderstormSampleContext,
		},
		SampleId: sampleId,
	}
}

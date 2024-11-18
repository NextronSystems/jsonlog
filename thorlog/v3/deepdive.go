package thorlog

import (
	"fmt"

	"github.com/NextronSystems/jsonlog"
)

type DeepDiveChunk struct {
	jsonlog.ObjectHeader

	Target *File `json:"file" textlog:"file"`

	ChunkOffset HexNumber   `json:"chunk_offset" textlog:"chunk_offset"`
	ChunkEnd    HexNumber   `json:"chunk_end" textlog:"chunk_end"`
	Content     *SparseData `json:"content" textlog:"content,expand"`

	// BeaconConfig contains information about the beacon configuration if the chunk contains a CS beacon.
	BeaconConfig *BeaconConfig `json:"beacon_config,omitempty" textlog:"beacon,expand,omitempty"`
}

func (DeepDiveChunk) reportable() {}

type HexNumber uint64

func (h HexNumber) String() string {
	return fmt.Sprintf("%#x", uint64(h))
}

const typeDeepDiveChunk = "file chunk"

func init() { AddLogObjectType(typeDeepDiveChunk, &DeepDiveChunk{}) }

func NewDeepDiveChunk() *DeepDiveChunk {
	return &DeepDiveChunk{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeDeepDiveChunk,
		},
	}
}

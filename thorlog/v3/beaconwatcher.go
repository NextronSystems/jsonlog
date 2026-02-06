package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

const typeNetworkConnectingThread = "network connecting thread"

func init() { AddLogObjectType(typeNetworkConnectingThread, &NetworkConnectingThread{}) }

type NetworkConnectingThread struct {
	jsonlog.ObjectHeader

	ThreadId uint32   `json:"thread_id" textlog:"thread_id"`
	Process  *Process `json:"process" textlog:",expand"`

	CallbackInterval time.Duration      `json:"callback_interval" textlog:"callback_interval"`
	Connections      NetworkConnections `json:"connections" textlog:"connections"`
}

func (NetworkConnectingThread) observed() {}

func NewNetworkConnectingThread(threadId uint32, process *Process) *NetworkConnectingThread {
	return &NetworkConnectingThread{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeNetworkConnectingThread,
		},
		ThreadId: threadId,
		Process:  process,
	}
}

type NetworkConnection struct {
	Protocol string `json:"protocol"`
	Server   string `json:"server"`
}

func (n NetworkConnection) String() string {
	return n.Server + "(" + n.Protocol + ")"
}

type NetworkConnections []NetworkConnection

func (n NetworkConnections) String() string {
	var s string
	for i, c := range n {
		if i > 0 {
			s += ", "
		}
		s += c.String()
	}
	return s
}

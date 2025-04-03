package thorlog

import (
	"fmt"
	"strings"
	"time"

	"github.com/NextronSystems/jsonlog"
)

type Process struct {
	jsonlog.ObjectHeader

	Pid int32 `json:"pid" textlog:"pid"`

	Dead bool `json:"dead,omitempty" textlog:"dead,omitempty"`
	// Only filled if dead is false
	ProcessInfo `textlog:",expand,omitempty"`

	// BeaconConfig contains information about a Cobalt Strike Beacon if the process contains one.
	BeaconConfig *BeaconConfig `json:"beacon_config,omitempty" textlog:"beacon,expand,omitempty"`

	// PeSieveReport contains information from PE-Sieve about the process, if any exists.
	PeSieveReport *PeSieveReport `json:"pe_sieve,omitempty" textlog:"pe_sieve,expand,omitempty"`
}

func (Process) reportable() {}

type ProcessInfo struct {
	Name    string `json:"name" textlog:"name"`
	Cmdline string `json:"command" textlog:"command"`
	User    string `json:"owner" textlog:"owner"`

	Image *File `json:"image" textlog:"image,expand"`

	ParentInfo struct {
		Pid         int32  `json:"pid" textlog:"ppid"`
		Exe         string `json:"exe" textlog:"parent"`
		CommandLine string `json:"command" textlog:"parent_command"`
	} `json:"parent_info,omitempty" textlog:",expand,omitempty"`

	ProcessTree StringList `json:"tree" textlog:"tree,omitempty" jsonschema:"nullable"`

	Created Time   `json:"created" textlog:"created"`
	Session string `json:"session" textlog:"session,omitempty"`

	ProcessConnections `textlog:",expand"`

	Memory *SparseData `json:"memory,omitempty" textlog:"memory,expand,omitempty"`
}

type ProcessConnections struct {
	ListenPorts     ProcessListenPorts  `json:"listen_ports" textlog:"listen_ports,omitempty" jsonschema:"nullable"`
	Connections     []ProcessConnection `json:"connections" textlog:"-" jsonschema:"nullable"`
	ConnectionCount int                 `json:"-" textlog:"connection_count"`
}

type ProcessListenPorts []uint32

func (p ProcessListenPorts) String() string {
	var listenPortsStr []string
	for _, port := range p {
		listenPortsStr = append(listenPortsStr, fmt.Sprint(port))
	}
	return strings.Join(listenPortsStr, ",")
}

type ProcessConnection struct {
	Fd uint32 `json:"-" textlog:"-"`
	// Status is the connection status, e.g. ESTABLISHED, LISTEN, etc.
	Status     string `json:"status" textlog:"-"`
	Ip         string `json:"ip" textlog:"ip"`
	Port       uint32 `json:"port" textlog:"port"`
	RemoteIp   string `json:"rip,omitempty" textlog:"rip,omitempty"`
	RemotePort uint32 `json:"rport,omitempty" textlog:"rport,omitempty"`
	// Protocol is the layer 4 protocol used for the connection, e.g. TCP, UDP, etc.
	Protocol string `json:"protocol,omitempty" textlog:"protocol,omitempty"`
}

const typeProcess = "process"

func init() { AddLogObjectType(typeProcess, &Process{}) }

func NewProcess(pid int32) *Process {
	return &Process{
		ObjectHeader: jsonlog.ObjectHeader{
			Type:    typeProcess,
			Summary: fmt.Sprintf("PID %d", pid),
		},
		Pid: pid,
	}
}

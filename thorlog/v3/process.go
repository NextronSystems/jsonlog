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

	ParentInfo ParentProcessInfo `json:"parent_info,omitempty" textlog:",expand,omitempty"`

	ProcessTree StringList `json:"tree" textlog:"tree,omitempty" jsonschema:"nullable"`

	Created time.Time `json:"created" textlog:"created"`
	Session string    `json:"session" textlog:"session,omitempty"`

	ProcessConnections `textlog:",expand"`

	Sections Sections `json:"sections,omitempty" textlog:"-"`
}

type ParentProcessInfo struct {
	Pid         int32  `json:"pid" textlog:"ppid"`
	Exe         string `json:"exe" textlog:"parent"`
	CommandLine string `json:"command" textlog:"parent_command"`
}

type Sections []Section

// Section describes a memory range in a process's virtual memory.
// This typically corresponds to a section in an executable file or library, such as .text, .data, etc.,
// or a stack, heap, or similar.
// In Linux terms: it corresponds to a line in /proc/<pid>/maps.
type Section struct {
	// Name of the section. For sections from loaded libraries, this is the library's file path.
	// For other memory ranges, this is OS specific and may be empty.
	Name string `json:"name"`
	// Address is the start address of the section in the process's virtual memory.
	Address uint64 `json:"address"`
	// Size is the size of the section in bytes.
	Size uint64 `json:"size" textlog:"size"`
	// Offset is the offset within the mapped file or library, if this section
	// corresponds to a file section. If this section does not correspond to a file,
	// this is empty.
	Offset uint64 `json:"offset,omitempty"`
	// SparseData contains a sparse representation of the section's data.
	// Only the interesting parts of the section are included, typically those that have been matched.
	SparseData *SparseData `json:"sparse_data,omitempty"`
	// Permissions of the section.
	Permissions RwxPermissions `json:"permissions"`
}

// RelativeTextPointer implements the jsonlog.TextReferenceResolver interface for Sections.
// It resolves a reference to a Section's SparseData field to a human-readable string.
func (s *Sections) RelativeTextPointer(pointee any) (string, bool) {
	for i := range *s {
		section := &(*s)[i]
		if pointee == &section.SparseData {
			if section.Name != "" {
				return section.Name, true
			} else {
				return fmt.Sprintf("0x%x", section.Address), true
			}
		}
	}
	return "", false
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
	RemoteIp   string `json:"remote_ip,omitempty" textlog:"rip,omitempty"`
	RemotePort uint32 `json:"remote_port,omitempty" textlog:"rport,omitempty"`
	// Protocol is the layer 4 protocol used for the connection, e.g. TCP, UDP, etc.
	Protocol string `json:"protocol,omitempty" textlog:"protocol,omitempty"`
}

const typeProcess = "process"

func init() { AddLogObjectType(typeProcess, &Process{}) }

func NewProcess(pid int32) *Process {
	return &Process{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeProcess,
		},
		Pid: pid,
	}
}

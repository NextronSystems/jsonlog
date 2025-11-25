package thorlog

import (
	"time"
)

// EBPFProgram describes an eBPF program attached to a specific endpoint in the kernel.
//
// To use eBPF nomenclature: This struct describes an eBPF link and its corresponding program.
// The exposed information by the kernel about links can be found at
// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/include/uapi/linux/bpf.h?h=v6.17#n6680,
// and program information at
// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/include/uapi/linux/bpf.h?h=v6.17#n6610.
//
// eBPF programs can be attached to a wide range of things; the LinkType contains what sort of object
// the program is attached to, and AttachTarget contains what specific object it is attached to.
//
// EBPFProgram has a content that contains the (kernel translated) instructions,
// provided that the kernel does not hide them due to the kernel.kptr_restrict sysctl.
type EBPFProgram struct {
	LogObjectHeader

	// Tag is a hash calculated by the kernel over the program instructions.
	// It can be used to uniquely identify the attached program.
	Tag string `textlog:"tag" json:"tag"`
	// User that loaded the eBPF program
	User string `textlog:"user" json:"user"`
	// Program name
	Name string `textlog:"name" json:"name"`
	// Size of the loaded program.
	//
	// This relates to instructions that have already been rewritten by the kernel;
	// as such, it does not have to be the exact size of the instructions that were passed
	// when the program was loaded.
	Size uint64 `textlog:"size" json:"size"`
	// Maps used by this program
	Maps []string `json:"maps"`
	// Functions declared by this program
	Functions []string `json:"functions"`
	// Timestamp when this program was loaded
	LoadTime time.Time `textlog:"load_time" json:"load_time"`
	// RAM locked by this eBPF program
	MemoryLocked uint64 `json:"memory_locked"`
	// Type of object the eBPF program is attached to (kprobe, syscall, tracepoint, ...)
	LinkType string `textlog:"link_type" json:"link_type"`
	// eBPF program type, i.e. whether this is a program for packet inspection / kprobe / tracepoint / ...
	ProgramType string `json:"program_type"`
	// The object the eBPF program is attached to.
	//
	// Depending on the LinkType, different fields will be present in this struct.
	AttachTarget EBPFAttachTarget `textlog:",expand" json:"attach_target"`
	// Content contains extracts from the kernel translated instructions that are
	// relevant for matches on this program.
	Content *SparseData `json:"content,omitempty"`
}

// EBPFAttachTarget describes the target that a BPF program is attached to.
type EBPFAttachTarget struct {
	// uprobe / tracepoint / cgroup specific; the path of the hooked ELF / tracepoint / cgroup, respectively
	Path string `textlog:"path,omitempty" json:"path,omitempty"`
	// uprobe specific; the PID of the hooked process, or nothing if the probe is for all processes
	Pid uint32 `textlog:"pid,omitempty" json:"pid,omitempty"`
	// uprobe / kprobe specific; the symbols that are hooked
	Symbols StringList `textlog:"symbol,omitempty" json:"symbols,omitempty"`
	// netkit / TCX / XDP specific; Network interface that the eBPF is attached to
	Interface string `textlog:"interface,omitempty" json:"interface,omitempty"`
	// netns / tracing / perf event specific; ID of the object attached to
	ObjectId int64 `textlog:"object_id,omitempty" json:"object_id,omitempty"`
	// netfilter specific; Protocol family (IPv4 or IPv6)
	Protocol string `textlog:"protocol,omitempty" json:"protocol,omitempty"`
	// netfilter specific; Hook (prerouting, postrouting, forward, local in, or local out)
	Hook string `textlog:"hook,omitempty" json:"hook,omitempty"`
	// netfilter specific; Priority (lower is executed earlier)
	Priority int `textlog:"priority,omitempty" json:"priority,omitempty"`
}

func (EBPFProgram) reportable() {}

const typeEbpfProgram = "eBPF program"

func init() { AddLogObjectType(typeEbpfProgram, &EBPFProgram{}) }

func NewEBPFProgram() *EBPFProgram {
	return &EBPFProgram{
		LogObjectHeader: LogObjectHeader{
			Type: typeEbpfProgram,
		},
	}
}

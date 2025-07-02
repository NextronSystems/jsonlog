package thorlog

type LinuxKernelModule struct {
	LogObjectHeader

	Name string `json:"name" textlog:"name"`
	Size int    `json:"size,omitempty" textlog:"size,omitempty"`

	// Whether this modules was compiled into the kernel
	IncludedInKernel bool `json:"included_in_kernel" textlog:"included_in_kernel"`

	Refcount   int          `json:"ref_count"`
	UsedBy     StringList   `json:"used_by"`
	Version    string       `json:"version"`
	Parameters KeyValueList `json:"parameters,omitempty" textlog:"parameters,omitempty"`

	File        *File      `json:"file" textlog:"file,expand,omitempty"`
	Description StringList `json:"description" textlog:"description"`
	Author      string     `json:"author" textlog:"author"`
}

func (LinuxKernelModule) reportable() {}

const typeLinuxKernelModule = "Linux kernel module"

func NewLinuxKernelModule(name string) *LinuxKernelModule {
	return &LinuxKernelModule{
		LogObjectHeader: LogObjectHeader{
			Type: typeLinuxKernelModule,
		},
		Name: name,
	}
}

func init() {
	AddLogObjectType(typeLinuxKernelModule, &LinuxKernelModule{})
}

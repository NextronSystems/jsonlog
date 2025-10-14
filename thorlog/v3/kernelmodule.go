package thorlog

type LinuxKernelModule struct {
	LogObjectHeader

	Name string `json:"name" textlog:"name"`
	Size int    `json:"size,omitempty" textlog:"size,omitempty"`

	// Whether this modules was compiled into the kernel
	IncludedInKernel bool `json:"included_in_kernel" textlog:"included_in_kernel"`

	Refcount int        `json:"ref_count"`
	UsedBy   StringList `json:"used_by"`
	// List of modules that this module depends on (from /proc/modules)
	DependsOn  StringList   `json:"depends_on,omitempty"`
	Version    string       `json:"version"`
	Parameters KeyValueList `json:"parameters,omitempty" textlog:"parameters,omitempty"`
	// Current load state of the module: "Live", "Loading", or "Unloading" (from /proc/modules)
	LoadState string `json:"load_state,omitempty"`

	File        *File      `json:"file" textlog:"file,expand,omitempty"`
	Description StringList `json:"description" textlog:"description"`
	Author      string     `json:"author" textlog:"author"`

	// Indicates if this module was found in /proc/modules (currently loaded modules)
	InProcModules bool `json:"in_proc_modules" textlog:"in_proc_modules"`

	// Indicates the kernel exposes this module under /sys/module (sysfs entry present).
	InSysModule bool `json:"in_sys_module" textlog:"in_sys_module"`

	// Indicates if this module was found in /proc/vmallocinfo (modules with vmalloc allocations)
	InVmallocinfo bool `json:"in_vmallocinfo" textlog:"in_vmallocinfo"`
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

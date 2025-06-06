package thorlog

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NextronSystems/jsonlog"
)

type HostInfo struct {
	jsonlog.ObjectHeader

	Hostname    string          `json:"hostname" textlog:"hostname"`
	Domain      string          `json:"domain" textlog:"domain,omitempty"`
	Platform    PlatformInfo    `json:"platform" textlog:",expand"`
	Uptime      time.Duration   `json:"uptime" textlog:"uptime"`
	Cpus        int             `json:"cpu_count" textlog:"cpu_count"`
	Memory      uint64          `json:"memory" textlog:"memory"`
	Timezone    string          `json:"timezone" textlog:"timezone"`
	Language    string          `json:"language" textlog:"language"`
	Interfaces  []InterfaceInfo `json:"interfaces" textlog:",expand"`
	SystemType  SystemType      `json:"system_type" textlog:"system_type"`
	MountPoints []MountInfo     `json:"mount_points"`
}

const typeHostInfo = "system information"

func init() { AddLogObjectType(typeHostInfo, &HostInfo{}) }

func NewHostInfo() *HostInfo {
	return &HostInfo{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeHostInfo,
		},
	}
}

func (h *HostInfo) UnmarshalJSON(data []byte) error {
	// PlatformInfo is an embedded object, so we need to unmarshal first into a struct
	// that correctly resolves it based on the type.
	// To do this, we create a struct that has the same fields as HostInfo,
	// but with the Platform field as an EmbeddedObject and without the UnmarshalJSON method.
	type hostInfoClone HostInfo
	var unmarshalableInfo struct {
		hostInfoClone
		Platform EmbeddedObject `json:"platform"`
	}
	err := json.Unmarshal(data, &unmarshalableInfo)
	if err != nil {
		return err
	}
	*h = HostInfo(unmarshalableInfo.hostInfoClone)
	if platformInfo, isPlatformInfo := unmarshalableInfo.Platform.Object.(PlatformInfo); isPlatformInfo {
		h.Platform = platformInfo
	} else {
		return fmt.Errorf("platform information has invalid type %s", unmarshalableInfo.Platform.Object.EmbeddedHeader().Type)
	}
	return nil
}

type SystemType string

const (
	SystemTypeServer           SystemType = "Server"
	SystemTypeWorkstation      SystemType = "Workstation"
	SystemTypeDomainController SystemType = "Domain Controller"
)

type MountInfo struct {
	// FSType is the filesystem that is mounted, e.g. ext4, ntfs, etc.
	FSType string `json:"fs_type"`
	// Source is the OS description of the source of the mount.
	// This can differ greatly between OSes and filesystems.
	// For example, on Linux, for local partitions, this is the device path.
	Source string `json:"source"`
	// Target is the path where the filesystem is mounted.
	Target string `json:"target"`
	// Class is the class of the mount, e.g. local, network, removable, etc.
	// This determines how the mount is treated by THOR.
	// It is not innately part of the mount information, but is determined by THOR.
	Class string `json:"class"`
}

type InterfaceInfo struct {
	Name        string `json:"name"`
	IpAddress   string `json:"ip_address" textlog:"ip_address"`
	Ipv6Address string `json:"ipv6_address,omitempty"`
	MacAddress  string `json:"mac_address,omitempty"`
}

type PlatformInfo interface {
	jsonlog.Object
	platform()
}

type PlatformInfoMacos struct {
	jsonlog.ObjectHeader

	Name          string `json:"name" textlog:"name"`
	Version       string `json:"version" textlog:"version"`
	KernelName    string `json:"kernel_name" textlog:"kernel_name"`
	KernelVersion string `json:"kernel_version" textlog:"kernel_version"`
	Proc          string `json:"proc" textlog:"proc"`
	Arch          string `json:"arch" textlog:"arch"`
}

func (PlatformInfoMacos) platform() {}

const typePlatformInfoMacos = "MacOS platform information"

func init() { AddLogObjectType(typePlatformInfoMacos, &PlatformInfoMacos{}) }

func NewMacOSPlatformInfo() *PlatformInfoMacos {
	return &PlatformInfoMacos{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typePlatformInfoMacos,
		},
	}
}

type PlatformInfoLinux struct {
	jsonlog.ObjectHeader

	Name          string `json:"name" textlog:"name"`
	KernelName    string `json:"kernel_name" textlog:"kernel_name"`
	KernelVersion string `json:"kernel_version" textlog:"kernel_version"`
	Proc          string `json:"proc" textlog:"proc"`
	Arch          string `json:"arch" textlog:"arch"`
}

func (PlatformInfoLinux) platform() {}

const typePlatformInfoLinux = "Linux platform information"

func init() { AddLogObjectType(typePlatformInfoLinux, &PlatformInfoLinux{}) }

func NewLinuxPlatformInfo() *PlatformInfoLinux {
	return &PlatformInfoLinux{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typePlatformInfoLinux,
		},
	}
}

type PlatformInfoWindows struct {
	jsonlog.ObjectHeader

	Name        string    `json:"name" textlog:"name"`
	Type        string    `json:"type" textlog:"type"`
	Version     string    `json:"version" textlog:"version"`
	Proc        string    `json:"proc" textlog:"proc"`
	Arch        string    `json:"arch" textlog:"arch"`
	InstalledOn time.Time `json:"installed_on" textlog:"installed_on"`
	BuildNumber string    `json:"build_number" textlog:"build_number"`
}

func (PlatformInfoWindows) platform() {}

const typePlatformInfoWindows = "Windows platform information"

func init() { AddLogObjectType(typePlatformInfoWindows, &PlatformInfoWindows{}) }

func NewWindowsPlatformInfo() *PlatformInfoWindows {
	return &PlatformInfoWindows{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typePlatformInfoWindows,
		},
	}
}

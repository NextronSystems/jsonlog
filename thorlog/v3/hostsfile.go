package thorlog

type HostsFileEntry struct {
	LogObjectHeader
	Host string `json:"host" textlog:"host"`
	IP   string `json:"ip" textlog:"ip"`
}

func (HostsFileEntry) observed() {}

const typeHostsFileEntry = "hosts file entry"

func init() { AddLogObjectType(typeHostsFileEntry, &HostsFileEntry{}) }

func NewHostsFileEntry(host string, ip string) *HostsFileEntry {
	return &HostsFileEntry{
		LogObjectHeader: LogObjectHeader{
			Type: typeHostsFileEntry,
		},
		Host: host,
		IP:   ip,
	}
}

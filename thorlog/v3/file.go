package thorlog

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NextronSystems/jsonlog"
)

type File struct {
	jsonlog.ObjectHeader

	// Path is the full path of the file (possibly including archives, e.g. /path/to/archive.zip/file.txt)
	Path string `json:"path" textlog:"file"`

	// Exists is a flag indicating whether the file exists or not. This is useful for files that are referenced elsewhere, but do not necessarily exist.
	Exists Existence `json:"exists" textlog:"exists,omitempty"`

	// Extension is the file extension of the file (e.g. .txt, .exe, etc.)
	Extension string `json:"extension" textlog:"extension,omitempty"`

	// FileMode is the type of the file (e.g. file, directory, symlink, etc.)
	FileMode FileModeType `json:"-" textlog:"-"`

	// MagicHeader is the magic header of the file (e.g. PE, ZIP, etc.)
	MagicHeader string `json:"magic_header" textlog:"type,omitempty"`

	// FileHashes contains the MD5, SHA1, and SHA256 hashes of the file, provided that the file is regular and could be read
	Hashes *FileHashes `json:"hashes,omitempty" textlog:",expand,omitempty"`

	// FirstBytes contains the first bytes of the file
	FirstBytes FirstBytes `json:"first_bytes,omitempty" textlog:"firstbytes,omitempty"`

	// Filetimes contains the file times of the file (e.g. created, modified, accessed, etc.)
	Filetimes *Filetimes `json:"file_times,omitempty" textlog:",expand,omitempty"`

	Size uint64 `json:"size" textlog:"size,omitempty"`

	// Permissions contains the permissions of the file. This can be either Unix or Windows permissions.
	Permissions Permissions `json:"permissions" textlog:",expand,omitempty"`

	// PeInfo contains information about the PE file, if the file is a PE file
	PeInfo *PeInfo `json:"pe_info,omitempty" textlog:",expand,omitempty"`

	// Target is only set for symlinks and contains the target path of the symlink
	Target string `json:"target,omitempty" textlog:"target,omitempty"`

	// UnpackSource is set for files that originate from another, unpacked file (possibly with multiple layers of unpacking)
	UnpackSource ArrowStringList `json:"unpack_source,omitempty" textlog:"unpack_source,omitempty" jsonschema:"nullable"`

	// LinkInfo contains information about the link, if the file is a windows link file (.lnk)
	LinkInfo *LinkInfo `json:"link_info,omitempty" textlog:",expand,omitempty"`

	// RecycleBinInfo contains information about the file if it was in the recycle bin
	RecycleBinInfo *RecycleBinIndexFile `json:"recycle_bin_info,omitempty" textlog:",expand,omitempty"`

	// WerInfo contains information about the file if it was a Windows Error Reporting crash report
	WerInfo *WerCrashReport `json:"wer_info,omitempty" textlog:",expand,omitempty"`

	// Content contains extracts from the content of the file, typically focusing on any matched patterns.
	Content *SparseData `json:"content,omitempty" textlog:"content,expand,omitempty"`

	// BeaconConfig contains information about a Cobalt Strike Beacon if the file contains one.
	BeaconConfig *BeaconConfig `json:"beacon_config,omitempty" textlog:"beacon,expand,omitempty"`

	// VirusTotalInfo contains information about the file from VirusTotal
	VirusTotalInfo *VirusTotalInformation `json:"virustotal,omitempty" textlog:"virustotal,expand,omitempty"`
}

func (File) reportable() {}

func (f *File) UnmarshalJSON(data []byte) error {
	// Permissions are either unix or windows permissions, so we need to try both
	type plainFile File

	var testFile struct {
		plainFile
		Permissions EmbeddedObject `json:"permissions"`
	}
	err := json.Unmarshal(data, &testFile)
	if err != nil {
		return err
	}
	perms, isPermissions := testFile.Permissions.Object.(Permissions)
	if !isPermissions && testFile.Permissions.Object != nil {
		return fmt.Errorf("invalid permissions type: %T", testFile.Permissions.Object)
	}
	*f = File(testFile.plainFile)
	f.Permissions = perms
	return nil
}

type FileHashes struct {
	Md5    string `json:"md5" textlog:"md5"`
	Sha1   string `json:"sha1" textlog:"sha1"`
	Sha256 string `json:"sha256" textlog:"sha256"`
}

type RecycleBinIndexFile struct {
	Version          uint64    `json:"-" textlog:"-"`
	OriginalFilename string    `json:"original_file_name" textlog:"original_filename"`
	DeletionTime     time.Time `json:"deletion_time" textlog:"deletion_time"`
	OriginalFilesize uint64    `json:"original_file_size" textlog:"-"`
}

type LinkInfo struct {
	Target       string    `json:"target" textlog:"target"`
	Arguments    string    `json:"arguments" textlog:"arguments"`
	CommandLine  string    `json:"command_line" textlog:"command_line"`
	CreationTime time.Time `json:"created" textlog:"-"`
	WriteTime    time.Time `json:"modified" textlog:"-"`
	AccessTime   time.Time `json:"accessed" textlog:"-"`
}

const ModifierWithMilliseconds = "with_millis"

type Filetimes struct {
	Mtime time.Time  `json:"modified" textlog:"modified,with_millis"`
	Atime *time.Time `json:"accessed,omitempty" textlog:"accessed,omitempty,with_millis"`
	Ctime *time.Time `json:"changed,omitempty" textlog:"changed,omitempty,with_millis"`
	Btime *time.Time `json:"created,omitempty" textlog:"created,omitempty,with_millis"`

	// Timestamps that are not always available, but only set if timestomping is detected
	UsnChangeTime       *time.Time `json:"usn_change_time,omitempty" textlog:"usn_change_time,omitempty,with_millis"`
	MftFileNameModified *time.Time `json:"mft_file_name_modified,omitempty" textlog:"mft_file_name_modified,omitempty,with_millis"`
	MftFileNameAccessed *time.Time `json:"mft_file_name_accessed,omitempty" textlog:"mft_file_name_accessed,omitempty,with_millis"`
	MftFileNameChanged  *time.Time `json:"mft_file_name_changed,omitempty" textlog:"mft_file_name_changed,omitempty,with_millis"`
	MftFileNameCreated  *time.Time `json:"mft_file_name_created,omitempty" textlog:"mft_file_name_created,omitempty,with_millis"`
}

type PeInfo struct {
	Company         string `json:"company" textlog:"company,omitempty"`
	FileDescription string `json:"description" textlog:"description,omitempty"`
	LegalCopyright  string `json:"legal_copyright" textlog:"legal_copyright,omitempty"`
	Product         string `json:"product" textlog:"product,omitempty"`
	OriginalName    string `json:"original_name" textlog:"original_name,omitempty"`
	InternalName    string `json:"internal_name" textlog:"internal_name,omitempty"`

	Signed     bool            `json:"signed" textlog:"signed"`
	Signatures []SignatureInfo `json:"signatures" textlog:",expand" jsonschema:"nullable"`

	Imphash           string    `json:"imphash" textlog:"imphash,omitempty"`
	RichHeaderHash    string    `json:"rich_header_hash"`
	CreationTimestamp time.Time `json:"creation_timestamp"`
}

type SignatureInfo struct {
	CertificateName string `json:"certificate_name" textlog:"certificate_name,omitempty"`
	SignatureValid  bool   `json:"signature_valid" textlog:"signature_valid"`
}

type FileModeType string

const (
	Undefined FileModeType = "undefined"
	Directory FileModeType = "directory"
	Irregular FileModeType = "irregular"
	Symlink   FileModeType = "symlink"
	ModeFile  FileModeType = "file"
)

type Existence string

const (
	ExistenceYes                 Existence = "yes"
	ExistenceNo                  Existence = "no"
	ExistenceUnknown             Existence = "unknown"
	ExistenceDisappeared         Existence = "disappeared"          // Unknown because disappeared
	ExistenceExpansionInfeasible Existence = "expansion_infeasible" // Unknown because expansion
	ExistenceNonLocal            Existence = "nonlocal"             // Unknown because not local
	ExistenceExcluded            Existence = "excluded"             // Unknown because excluded
)

func (e Existence) IsZero() bool {
	return e == ExistenceYes
}

const typeFile = "file"

func init() { AddLogObjectType(typeFile, &File{}) }

func NewFile(path string) *File {
	return &File{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeFile,
		},
		Path: path,
	}
}

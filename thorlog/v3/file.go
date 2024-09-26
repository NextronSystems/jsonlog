package thorlog

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/NextronSystems/jsonlog"
)

type File struct {
	jsonlog.ObjectHeader
	Path      string    `json:"path" textlog:"file"`
	Exists    Existence `json:"exists" textlog:"exists,omitempty"`
	Extension string    `json:"extension" textlog:"extension,omitempty"`

	FileMode FileModeType `json:"-" textlog:"-"`

	MagicHeader string `json:"magic_header" textlog:"type,omitempty"`

	Hashes     *FileHashes `json:"hashes,omitempty" textlog:",expand,omitempty"`
	FirstBytes FirstBytes  `json:"firstbytes,omitempty" textlog:"firstbytes,omitempty"`

	Filetimes *Filetimes `json:"filetimes,omitempty" textlog:",expand,omitempty"`

	Size uint64 `json:"size" textlog:"size,omitempty"`

	Permissions Permissions `json:"permissions" textlog:",expand,omitempty"`

	PeInfo *PeInfo `json:"pe_info,omitempty" textlog:",expand,omitempty"`

	// Target is only set for symlinks and contains the target path of the symlink
	Target string `json:"target,omitempty" textlog:"target,omitempty"`

	// UnpackSource is set for files that originate from another, unpacked file (possibly with multiple layers of unpacking)
	UnpackSource ArrowStringList `json:"unpack_source,omitempty" textlog:"unpack_source,omitempty" jsonschema:"nullable"`

	LinkInfo *LinkInfo `json:"link_info,omitempty" textlog:",expand,omitempty"`

	RecycleBinInfo *RecycleBinIndexFile `json:"recycle_bin_info,omitempty" textlog:",expand,omitempty"`

	WerInfo *WerCrashReport `json:"wer_info,omitempty" textlog:",expand,omitempty"`

	Content *SparseData `json:"content,omitempty" textlog:"content,expand,omitempty"`
}

func (f *File) UnmarshalJSON(data []byte) error {
	// Permissions are either unix or windows permissions, so we need to try both
	type plainFile File

	var testFile plainFile
	testFile.Permissions = &UnixPermissions{}
	err := json.Unmarshal(data, &testFile)
	if err != nil {
		testFile.Permissions = &WindowsPermissions{}
		err = json.Unmarshal(data, &testFile)
		if err != nil {
			return err
		}
	}
	*f = File(testFile)
	return nil
}

type FileHashes struct {
	Md5    string `json:"md5" textlog:"md5"`
	Sha1   string `json:"sha1" textlog:"sha1"`
	Sha256 string `json:"sha256" textlog:"sha256"`
}

type RecycleBinIndexFile struct {
	Version          uint64    `json:"-" textlog:"-"`
	OriginalFilename string    `json:"original_filename" textlog:"original_filename"`
	DeletionTime     time.Time `json:"deletion_time" textlog:"deletion_time"`
	OriginalFilesize uint64    `json:"original_filesize" textlog:"-"`
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

type FirstBytes []byte

func trimNonPrintableChars(data []byte) string {
	var builder strings.Builder
	builder.Grow(len(data))
	for _, char := range data {
		if char >= 0x20 && char <= 0x7E {
			builder.WriteByte(char)
		}
	}
	return builder.String()
}

func (f FirstBytes) String() string {
	return hex.EncodeToString(f) + " / " + trimNonPrintableChars(f)
}

type firstBytesJson struct {
	Hex   string `json:"hex"`
	Ascii string `json:"ascii"`
}

func (f FirstBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(firstBytesJson{hex.EncodeToString(f), trimNonPrintableChars(f)})
}

func (f *FirstBytes) UnmarshalJSON(data []byte) error {
	var jsonStruct firstBytesJson
	err := json.Unmarshal(data, &jsonStruct)
	if err != nil {
		return err
	}
	unhexedData, err := hex.DecodeString(jsonStruct.Hex)
	if err != nil {
		return err
	}
	*f = unhexedData
	return nil
}

func (f FirstBytes) JSONSchemaAlias() any { return firstBytesJson{} }

type FileModeType string

const (
	Undefined   FileModeType = "undefined"
	NotExisting FileModeType = "nonexistent"
	Directory   FileModeType = "directory"
	Irregular   FileModeType = "irregular"
	Symlink     FileModeType = "symlink"
	ModeFile    FileModeType = "file"
)

func (f FileModeType) AsExistence() Existence {
	if f == Undefined {
		return ExistenceUnknown
	} else if f == NotExisting {
		return ExistenceNo
	} else {
		return ExistenceYes
	}
}

type Existence int

const (
	ExistenceYes Existence = iota
	ExistenceUnknown
	ExistenceNo
)

func (e Existence) String() string {
	switch e {
	case ExistenceYes:
		return "yes"
	case ExistenceUnknown:
		return "unknown"
	case ExistenceNo:
		return "no"
	default:
		panic("invalid existence")
	}
}

func (e Existence) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *Existence) UnmarshalJSON(data []byte) error {
	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err != nil {
		return err
	}
	switch stringValue {
	case "yes":
		*e = ExistenceYes
	case "unknown":
		*e = ExistenceUnknown
	case "no":
		*e = ExistenceNo
	default:
		return fmt.Errorf("unknown existence %s", stringValue)
	}
	return nil
}

func (e Existence) JSONSchemaAlias() any { return "" }

const typeFile = "file"

func init() { AddLogObjectType(typeFile, &File{}) }

func NewFile(path string) *File {
	return &File{
		ObjectHeader: jsonlog.ObjectHeader{
			Type:    typeFile,
			Summary: path,
		},
		Path: path,
	}
}

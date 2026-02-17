package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
	"github.com/google/uuid"
)

// UALEntry holds information about a single entry of a User Access Log (UAL)
// database. These databases are written by the User Access Logging service which
// aggregates client usage data by roles and products.
//
// Reference: https://learn.microsoft.com/en-us/windows-server/administration/user-access-logging/get-started-with-user-access-logging
//
// A UALEntry represents a single entry in the CLIENTS table, possibly enriched with
// role details in the ROLE_IDS table of an accompanying SystemIdentity.mdb file.
//
//		from Current.mdb or <GUID>.mdb:
//		Table: 6                     CLIENTS (10)
//		     Number of columns:      374
//		     Column  Identifier      Name            Type
//		     1       1               RoleGuid        GUID
//		     2       2               TenantId        GUID
//		     3       3               TotalAccesses   Integer 32-bit unsigned
//		     4       4               InsertDate      Date and time
//		     5       5               LastAccess      Date and time
//		     6       128             Address         Binary data
//		     7       256             AuthenticatedUserName   Large text
//		     8       257             ClientName      Large text
//		     9       258             Day1            Integer 16-bit unsigned
//		     10      259             Day2            Integer 16-bit unsigned
//		     11      260             Day3            Integer 16-bit unsigned
//		     ...
//
//	 from SystemIdentity.mdb:
//	 Table: 7                    ROLE_IDS (12)
//	     Number of columns:      3
//	     Column  Identifier      Name            Type
//	     1       1               RoleGuid        GUID
//	     2       256             ProductName     Large text
//	     3       257             RoleName        Large text
type UALEntry struct {
	jsonlog.ObjectHeader

	// AuthenticatedUserName is the user name on the client that accompanies the UAL
	// entries from installed roles and products, if applicable.
	AuthenticatedUserName string `json:"authenticated_user_name" textlog:"authenticated_user_name"`
	// Address is the IP address of a client device that is used to access a role or
	// service.
	Address string `json:"address" textlog:"address"`
	// TotalAccesses is the number of times a particular user accessed a role or service.
	TotalAccesses uint32 `json:"total_accesses" textlog:"total_accesses"`
	// RoleGuid is the UAL assigned or registered GUID that represents the server role or
	// installed product.
	RoleGuid uuid.UUID `json:"role_guid" textlog:"role_guid"`
	// RoleName is the name of the role, component, or subproduct that is providing UAL
	// data.
	RoleName string `json:"role_name,omitempty" textlog:"role_name,omitempty"`
	// ProductName is the name of the software parent product, such as Windows, that is
	// providing UAL data. The value can be a GUID or a human-readable string.
	ProductName string `json:"product_name,omitempty" textlog:"product_name,omitempty"`
	// TenantId is a unique GUID for a tenant client of an installed role or product that
	// accompanies the UAL data, if applicable.
	TenantId uuid.UUID `json:"tenant_id" textlog:"tenant_id"`
	// InsertDate is the date and time when an IP address was first used to access a role
	// or service.
	InsertDate time.Time `json:"insert_date" textlog:"insert_date"`
	// LastAccess is the date and time when an IP address was last used to access a role
	// or service.
	LastAccess time.Time `json:"last_access" textlog:"last_access"`
	// ClientName. Usually unset.
	ClientName string `json:"client_name,omitempty" textlog:"client_name,omitempty"`
	// AccessesByDay is a map of the number of accesses per day of the year.
	AccessesByDay map[int]uint16 `json:"accesses_by_day" textlog:"-"`
}

const typeUALEntry = "User Access Log Entry"

func init() { AddLogObjectType(typeUALEntry, &UALEntry{}) }

func NewUALEntry() *UALEntry {
	return &UALEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeUALEntry,
		},
	}
}

func (UALEntry) observed() {}

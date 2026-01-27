package domain

type RoleID string
type RoleName string

const (
	// Role IDs (UUIDs matching seed data)
	RoleIDUser       RoleID = "00000000-0000-0000-0000-000000000001"
	RoleIDAdmin      RoleID = "00000000-0000-0000-0000-000000000002"
	RoleIDSuperAdmin RoleID = "00000000-0000-0000-0000-000000000003"

	// Role Names
	RoleNameUser       RoleName = "user"
	RoleNameAdmin      RoleName = "admin"
	RoleNameSuperAdmin RoleName = "super_admin"
)

// RoleIDToName maps Role IDs to their string names
var RoleIDToName = map[RoleID]RoleName{
	RoleIDUser:       RoleNameUser,
	RoleIDAdmin:      RoleNameAdmin,
	RoleIDSuperAdmin: RoleNameSuperAdmin,
}

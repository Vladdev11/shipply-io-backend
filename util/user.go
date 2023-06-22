package util

var UserRoleMap = map[int]string{
	1: OrganizationAdmin,
	2: OrganizationUser,
	3: ClientAdmin,
	4: ClientUser,
}

var (
	OrganizationAdmin    = "organization_admin"
	OrganizationAdminInt = 1
	OrganizationUser     = "organization_user"
	OrganizationUserInt  = 2
	ClientAdmin          = "client_admin"
	ClientAdminInt       = 3
	ClientUser           = "client_user"
	ClientUserInt        = 4
)

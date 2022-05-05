package entities

type RoleID int32

const (
	UserRoleID  RoleID = 1
	AdminRoleID RoleID = 2
)

type Role struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}

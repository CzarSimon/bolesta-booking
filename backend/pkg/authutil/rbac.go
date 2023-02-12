package authutil

const UserRole = "USER"

type Permission int

const (
	ReadBooking Permission = iota
	CreateBooking
	DeleteBooking
	ReadCabin
	ReadUser
)

type role struct {
	permissions map[Permission]bool
}

func (r role) Allowed(p Permission) bool {
	if r.permissions == nil {
		return false
	}

	_, ok := r.permissions[p]
	return ok
}

var userRole = role{
	permissions: map[Permission]bool{
		ReadBooking:   true,
		CreateBooking: true,
		DeleteBooking: true,
		ReadCabin:     true,
		ReadUser:      true,
	},
}

func getRole(roleName string) (role, bool) {
	if roleName == UserRole {
		return userRole, true
	}

	return role{}, false
}

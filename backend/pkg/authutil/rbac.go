package authutil

const AdminRole = "ADMIN"
const UserRole = "USER"
const AnonymousRole = "ANONYMOUS"

type Permission int

func (p Permission) String() string {
	switch p {
	case ReadBooking:
		return "ReadBooking"
	case CreateBooking:
		return "CreateBooking"
	case DeleteBooking:
		return "DeleteBooking"
	case ReadCabin:
		return "ReadCabin"
	case ReadUser:
		return "ReadUser"
	case CreateUser:
		return "CreateUser"
	default:
		return "Unkown"
	}
}

const (
	ReadBooking Permission = iota
	CreateBooking
	DeleteBooking
	ReadCabin
	ReadUser
	CreateUser
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

var adminRole = role{
	permissions: map[Permission]bool{
		ReadBooking: true,
		ReadCabin:   true,
		ReadUser:    true,
		CreateUser:  true,
	},
}

var anonymousRole = role{
	permissions: map[Permission]bool{},
}

func getRole(roleName string) (role, bool) {
	if roleName == UserRole {
		return userRole, true
	}

	if roleName == AdminRole {
		return adminRole, true
	}

	if roleName == AnonymousRole {
		return anonymousRole, true
	}

	return role{}, false
}

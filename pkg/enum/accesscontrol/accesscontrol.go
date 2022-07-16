package accesscontrol

import "errors"

type AccessControl string

const (
	MASTER  AccessControl = "MASTER"
	SYSTEM  AccessControl = "SYSTEM"
	ADMIN   AccessControl = "ADMIN"
	MEMBER  AccessControl = "MEMBER"
	PENDING AccessControl = "PENDING"
	BANNED  AccessControl = "BANNED"
)

var mapac map[string]AccessControl = map[string]AccessControl{
	"MASTER":  MASTER,
	"SYSTEM":  SYSTEM,
	"ADMIN":   ADMIN,
	"MEMBER":  MEMBER,
	"PENDING": PENDING,
	"BANNED":  BANNED,
}
var sliceac = [...]string{"MASTER", "SYSTEM", "ADMIN", "MEMBER", "PENDING", "BANNED"}

func (ac *AccessControl) Enum(str string) error {
	v, ok := mapac[str]
	if !ok {
		err := errors.New("accesscontrol: string doesn't match enum")
		return err
	}
	ac = &v
	return nil
}

func (ac *AccessControl) String() string {
	x := string(*ac)
	for _, v := range sliceac {
		if v == x {
			return x
		}
	}
	return ""
}

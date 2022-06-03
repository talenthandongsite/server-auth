package accesscontrol

import "errors"

type AccessControl string

const (
	MASTER  AccessControl = "MASTER"
	SYSTEM                = "SYSTEM"
	ADMIN                 = "ADMIN"
	MEMBER                = "MEMBER"
	PENDING               = "PENDING"
	BANNED                = "BANNED"
)

var mapper map[string]AccessControl = map[string]AccessControl{
	"MASTER":  MASTER,
	"SYSTEM":  SYSTEM,
	"ADMIN":   ADMIN,
	"MEMBER":  MEMBER,
	"PENDING": PENDING,
	"BANNED":  BANNED,
}

func Enum(str string) (AccessControl, error) {
	v, ok := mapper[str]
	if !ok {
		err := errors.New("accesscontrol: string doesn't match enum")
		return "", err
	}
	return v, nil
}

var accesscontrol = [...]string{"MASTER", "SYSTEM", "ADMIN", "MEMBER", "PENDING", "BANNED"}

func String(e AccessControl) string {
	x := string(e)
	for _, v := range accesscontrol {
		if v == x {
			return x
		}
	}
	return ""
}

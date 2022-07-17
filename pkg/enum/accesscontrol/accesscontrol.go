package accesscontrol

import (
	"errors"
	"strings"
)

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

func Enum(str string) (AccessControl, error) {
	v, ok := mapac[strings.ToUpper(str)]
	if !ok {
		err := errors.New("accesscontrol: string doesn't match enum")
		return "", err
	}
	return v, nil
}

func (ac *AccessControl) String() string {
	return string(*ac)
}

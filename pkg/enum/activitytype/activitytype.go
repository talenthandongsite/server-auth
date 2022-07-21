package activitytype

import (
	"errors"
	"strings"
)

type ActivityType string

const (
	CREATED         ActivityType = "CREATED"
	UPDATED         ActivityType = "UPDATED"
	ADMIN_NOTE      ActivityType = "ADMIN_NOTE"
	SIGN_IN         ActivityType = "SIGN_IN"
	KEYCHAIN_UPSERT ActivityType = "KEYCHAIN_UPSERT"
	KEYCHAIN_DELETE ActivityType = "KEYCHAIN_DELETE"
)

var mapat map[string]ActivityType = map[string]ActivityType{
	"CREATED":         CREATED,
	"UPDATED":         UPDATED,
	"ADMIN_NOTE":      ADMIN_NOTE,
	"SIGN_IN":         SIGN_IN,
	"KEYCHAIN_UPSERT": KEYCHAIN_UPSERT,
	"KEYCHAIN_DELETE": KEYCHAIN_DELETE,
}
var sliceat = [...]string{"CREATED", "UPDATED", "ADMIN_NOTE", "SIGN_IN", "KEYCHAIN_UPSERT", "KEYCHAIN_DELETE"}

func (at *ActivityType) Enum(str string) error {
	v, ok := mapat[strings.ToUpper(str)]
	if !ok {
		err := errors.New("accesscontrol: string doesn't match enum")
		return err
	}
	at = &v
	return nil
}

func (at *ActivityType) String() string {
	x := string(*at)
	for _, v := range sliceat {
		if v == x {
			return x
		}
	}

	return ""
}

package keychaintype

import (
	"errors"
	"strings"
)

type KeyChainType string

const (
	PASSWORD KeyChainType = "PASSWORD"
	KAKAO    KeyChainType = "KAKAO"
)

var mapkct map[string]KeyChainType = map[string]KeyChainType{
	"PASSWORD": PASSWORD,
	"KAKAO":    KAKAO,
}
var slicekct = [...]string{"PASSWORD", "KAKAO"}

func Enum(str string) (kct KeyChainType, err error) {
	v, ok := mapkct[strings.ToUpper(str)]
	if !ok {
		err := errors.New("accesscontrol: string doesn't match enum")
		return "", err
	}
	return v, nil
}

func (kct *KeyChainType) String() string {
	return string(*kct)
}

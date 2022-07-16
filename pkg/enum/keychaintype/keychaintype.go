package keychaintype

import "errors"

type KeyChainType string

const (
	PASSWORD KeyChainType = "PASSWORD"
	KAKAO                 = "KAKAO"
)

var mapkct map[string]KeyChainType = map[string]KeyChainType{
	"PASSWORD": PASSWORD,
	"KAKAO":    KAKAO,
}
var slicekct = [...]string{"PASSWORD", "KAKAO"}

func (kct *KeyChainType) Enum(str string) error {
	v, ok := mapkct[str]
	if !ok {
		err := errors.New("accesscontrol: string doesn't match enum")
		return err
	}
	kct = &v
	return nil
}

func (kct *KeyChainType) String() string {
	x := string(*kct)
	for _, v := range slicekct {
		if v == x {
			return x
		}
	}

	return ""
}

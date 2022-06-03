package keychaintype

type KeyChainType string

const (
	PASSWORD KeyChainType = "PASSWORD"
	KAKAO                 = "KAKAO"
)

var types = [...]string{"PASSWORD", "KAKAO"}

func String(e KeyChainType) string {

	x := string(e)
	for _, v := range types {
		if v == x {
			return x
		}
	}

	return ""
}

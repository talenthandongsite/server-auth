package keychaintype

import "testing"

func TestEnum(t *testing.T) {
	kct, err := Enum("password")
	if err != nil {
		t.Errorf("password is valid kct, got an error\n")
	}

	if kct != PASSWORD {
		t.Errorf("kct must be 'PASSWORD', got %s\n", kct)
	}
}

func TestString(t *testing.T) {
	kct := KAKAO
	str := kct.String()

	if str != "KAKAO" {
		t.Errorf("expected 'KAKAO', got %s", str)
	}
}

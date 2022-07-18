package accesscontrol

import "testing"

func TestEnum(t *testing.T) {
	ac, err := Enum("admin")
	if err != nil {
		t.Errorf("admin is valid ac, got an error\n")
	}

	if ac != ADMIN {
		t.Errorf("ac must be 'ADMIN', got %s\n", ac)
	}
}

func TestString(t *testing.T) {
	ac := MEMBER
	str := ac.String()

	if str != "MEMBER" {
		t.Errorf("expected 'MEMBER', got %s", str)
	}
}

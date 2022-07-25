package jwt

import (
	"testing"
	"time"
)

func TestForgeAndOpenToken(t *testing.T) {

	secret := []byte("1234")
	duration := time.Millisecond * 100000000
	config := JwtConfig{secret, duration}

	jwt, err := Init(config)
	if err != nil {
		t.Fatal(err)
	}

	id := "111"
	username := "1234"
	accessControl := "ADMIN"
	token, _, err := jwt.ForgeToken(id, username, accessControl)
	if err != nil {
		t.Fatal(err)
	}

	claim, err := jwt.OpenToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if claim.ID != id {
		t.Fatalf("ID is corrupted, expected %s, got %s\n", id, claim.ID)
	}
	if claim.Username != username {
		t.Fatalf("Username is corrupted, expected %s, got %s\n", username, claim.Username)
	}
	if claim.AccessControl != accessControl {
		t.Fatalf("ID is corrupted, expected %s, got %s\n", accessControl, claim.AccessControl)
	}
}

func TestDuration(t *testing.T) {
	secret := []byte("1234")
	duration := time.Millisecond * 100000000
	config := JwtConfig{secret, duration}

	jwt, err := Init(config)
	if err != nil {
		t.Fatal(err)
	}

	id := "111"
	username := "1234"
	accessControl := "ADMIN"
	token, _, err := jwt.ForgeToken(id, username, accessControl)
	if err != nil {
		t.Fatal(err)
	}

	claim, err := jwt.OpenToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if claim.ID != id {
		t.Fatalf("ID is corrupted, expected %s, got %s\n", id, claim.ID)
	}
	if claim.Username != username {
		t.Fatalf("Username is corrupted, expected %s, got %s\n", username, claim.Username)
	}
	if claim.AccessControl != accessControl {
		t.Fatalf("ID is corrupted, expected %s, got %s\n", accessControl, claim.AccessControl)
	}

	time.Sleep(time.Second * 3)

	_, err = jwt.OpenToken(token)
	t.Log(err)
	if err == nil {
		t.Fatal("Opening expired token should return err")
	}
}

package jwt

import (
	"context"
	"testing"
	"time"

	"github.com/talenthandongsite/server-auth/pkg/variable"
)

func TestForgeAndOpenToken(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, variable.TOKEN_SECRET, "1234")
	ctx = context.WithValue(ctx, variable.TOKEN_DURATION, "100000000")

	jwt, err := Init(ctx)
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
	ctx := context.Background()
	ctx = context.WithValue(ctx, variable.TOKEN_SECRET, "1234")
	ctx = context.WithValue(ctx, variable.TOKEN_DURATION, "2000")

	jwt, err := Init(ctx)
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

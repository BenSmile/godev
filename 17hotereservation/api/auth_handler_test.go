package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bensmile/hotel-reservation/db"
	"github.com/bensmile/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {

	user, err := types.NewUserFromParams(
		types.CreateUserParams{
			Password:  "pass",
			FirstName: "admin",
			LastName:  "admin",
			Email:     "admin@test.com"},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = userStore.InsertUser(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}

	return user

}

func TestAuthSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	insertedUser := insertTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/api/auth", authHandler.HandleLogin)

	params := types.AuthParams{
		Email:    "admin@test.com",
		Password: "pass",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status of 200 bug got %d", resp.StatusCode)
	}

	var authResp types.AuthResponse

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("expected the jwt token to be present in the auth response")
	}

	insertedUser.Password = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected the user to be the inserted user")
	}

}

func TestAuthWithWrongPasswordFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	insertTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/api/auth", authHandler.HandleLogin)

	params := types.AuthParams{
		Email:    "admin@test.com",
		Password: "passwrong",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status of 401 bug got %d", resp.StatusCode)
	}

	var authResp fiber.Map

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("expected the jwt token to be present in the auth response")
	}

	if authResp["message"] != "invalid credentials" {
		t.Fatalf("expected req response to have 'invalid credentials'")
	}

}

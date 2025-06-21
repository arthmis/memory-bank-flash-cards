package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"memorybank/database"
	"memorybank/queries"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
)

type User struct {
	UserId string `json:"user_id"`
}

type SessionToken struct {
	Object string `json:"object"`
	Jwt    string `json:"jwt"`
}

type SessionTokenLifetimeOverride struct {
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
}

func TestCreateDeck(t *testing.T) {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, "user=postgres dbname=postgres password=postgres port=7777")
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer dbpool.Close()

	deckQueries := queries.New(dbpool)

	env := &Env{
		Decks: database.DeckModel{DB: dbpool, Queries: *deckQueries},
	}

	godotenv.Load("../.env")
	clerkKey := os.Getenv("clerk_secret_key")

	client := resty.New()
	defer client.Close()
	id := "user_2ymUWK4LNISQ14RsaK8F5hPmzW0"

	session := clerk.Session{}
	res, err := client.R().SetBody(User{UserId: id}).SetHeader("Authorization", "Bearer "+clerkKey).SetHeader("Content-Type", "application/json").SetResult(&session).Post("https://api.clerk.com/v1/sessions")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, session.UserID == id, "True is true")
	assert.True(t, res.StatusCode() == 200, "True is true")

	sessionToken := SessionToken{}
	res, err = client.R().SetPathParam("session_id", session.ID).SetHeader("Authorization", "Bearer "+clerkKey).SetBody(SessionTokenLifetimeOverride{ExpiresInSeconds: 3600}).SetResult(&sessionToken).Post("https://api.clerk.com/v1/sessions/{session_id}/tokens")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, res.StatusCode() == 200, "True is true")

	handlers := NewStrictHandler(env, nil)
	e := echo.New()
	deck, _ := json.Marshal(CreateDeckJSONRequestBody{Name: "test"})
	req := httptest.NewRequest(http.MethodPost, "/api/decks", bytes.NewReader(deck))
	// req.Header.Set("Authorization", "Bearer "+sessionToken.Jwt)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	RegisterHandlers(e, handlers)

	if assert.NoError(t, handlers.CreateDeck(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	dbDeck, _ := deckQueries.GetDeck(ctx, "test")
	assert.Equal(t, dbDeck.Name, "test")
}

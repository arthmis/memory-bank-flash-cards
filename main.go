package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"memorybank/database"
	"memorybank/queries"
	"net/http"
	"os"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

type Env struct {
	decks database.DeckModel
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clerkKey := os.Getenv("clerk_secret_key")
	// Set the API key with your Clerk Secret Key
	clerk.SetKey(clerkKey)

	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, "user=postgres dbname=postgres password=postgres port=7777")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer dbpool.Close()

	deckQueries := queries.New(dbpool)

	env := &Env{
		decks: database.DeckModel{DB: dbpool, Queries: *deckQueries},
	}

	e := echo.New()
	e.GET("/dashboard", env.dashboard)
	e.POST("/deck", env.createDeck, Auth)
	e.Logger.Fatal(e.Start(":8000"))
}

func (env *Env) dashboard(c echo.Context) error {
	// handle getting all the decks and their names
	return c.NoContent(http.StatusOK)
}

// Process is the middleware function.
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the session JWT from the Authorization header
		sessionToken := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
		context := c.Request().Context()

		// Verify the session
		claims, err := jwt.Verify(context, &jwt.VerifyParams{
			Token: sessionToken,
		})
		if err != nil {
			// handle the error
			c.Response().WriteHeader(http.StatusUnauthorized)
			c.Response().Write([]byte(`{"access": "unauthorized"}`))
			return errors.New("unauthorized")
		}

		usr, err := user.Get(context, claims.Subject)
		if err != nil {
			// handle the error
		}
		fmt.Fprintf(c.Response().Writer, `{"user_id": "%s", "user_banned": "%t"}`, usr.ID, usr.Banned)

		next(c)

		return nil
	}
}

func (env *Env) createDeck(c echo.Context) error {
	name := c.FormValue("name")
	env.decks.CreateDeck(database.NewDeck{Name: name})
	return c.NoContent(http.StatusOK)
}

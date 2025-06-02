package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"memorybank/database"
	"memorybank/queries"
	"memorybank/views"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/clerk/clerk-sdk-go/v2"
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
	e.Static("/", "static")
	e.GET("/dashboard", env.dashboard)
	e.GET("/login", env.Login)
	// e.GET("api/login", env.login)
	e.POST("/deck", env.createDeck, Auth)
	e.Logger.Fatal(e.Start(":8000"))
}

func (env *Env) dashboard(c echo.Context) error {
	claims, ok := clerk.SessionClaimsFromContext(c.Request().Context())
	if !ok {
		c.Response().WriteHeader(http.StatusUnauthorized)
		c.Response().Write([]byte(`{"access": "unauthorized"}`))
		return errors.New("unauthorized")
	}
	fmt.Fprintf(c.Response().Writer, `{"user_id": "%s"}`, claims.Subject)

	// handle getting all the decks and their names
	return c.HTML(http.StatusOK, "hi")
}

func (env *Env) Login(c echo.Context) error {
	component := views.Login()
	return html(c, http.StatusOK, component)
}

func html(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

// Process is the middleware function.
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, ok := clerk.SessionClaimsFromContext(c.Request().Context())
		if !ok {
			c.Response().WriteHeader(http.StatusUnauthorized)
			c.Response().Write([]byte(`{"access": "unauthorized"}`))
			return errors.New("unauthorized")
		}
		fmt.Fprintf(c.Response().Writer, `{"user_id": "%s"}`, claims.Subject)

		next(c)

		return nil
	}
}

func (env *Env) createDeck(c echo.Context) error {
	name := c.FormValue("name")
	env.decks.CreateDeck(database.NewDeck{Name: name})
	return c.NoContent(http.StatusOK)
}

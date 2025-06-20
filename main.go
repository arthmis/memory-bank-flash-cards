package main

import (
	"context"
	"log"
	"memorybank/api"
	"memorybank/database"
	"memorybank/queries"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

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

	env := &api.Env{
		Decks: database.DeckModel{DB: dbpool, Queries: *deckQueries},
	}

	handlers := api.NewStrictHandler(env, nil)
	e := echo.New()
	api.RegisterHandlers(e, handlers)
	// e.Static("/dashboard", "app/dist")
	// // e.GET("/dashboard", env.dashboard, cookiesToAuth, handleAuth)
	// // protectedHandler := http.HandlerFunc(clerkAuth)
	// // headerAuthorization := clerkhttp.WithHeaderAuthorization()(protectedHandler)
	// // headerAuthorization := clerkhttp.WithHeaderAuthorization()
	// // authorization := echo.WrapMiddleware(headerAuthorization)
	// // e.GET("/dashboard", env.dashboard, cookiesToAuth, authorization, clerkAuth)
	e.Logger.Fatal(e.Start(":8000"))
}

// func (env *Env) dashboard(c echo.Context) error {
// 	component := views.Dashboard()
// 	return html(c, http.StatusOK, component)
// }

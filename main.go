package main

import (
	"context"
	"log"
	"memorybank/database"
	"memorybank/queries"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
)

type Env struct {
	decks database.DeckModel
}

func main() {
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
	e.POST("/deck", env.createDeck)
	e.Logger.Fatal(e.Start(":8000"))
}

func (env *Env) dashboard(c echo.Context) error {
	// handle getting all the decks and their names
	return c.NoContent(http.StatusOK)
}

func (env *Env) createDeck(c echo.Context) error {
	name := c.FormValue("name")
	env.decks.Queries.CreateDeck(c.Request().Context(), name)
	return c.NoContent(http.StatusOK)
}

package database

import (
	"context"
	"log"
	"memorybank/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DeckModel struct {
	DB      *pgxpool.Pool
	Queries queries.Queries
}

type NewDeck struct {
	name string
}

func (m DeckModel) save(newDeck NewDeck) {
	_, err := m.Queries.CreateDeck(context.Background(), newDeck.name)
	if err != nil {
		logger := log.Default()
		logger.Println(err)
	}
}

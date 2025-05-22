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
	Name string
}

func (m DeckModel) CreateDeck(newDeck NewDeck) {
	_, err := m.Queries.CreateDeck(context.Background(), newDeck.Name)
	if err != nil {
		logger := log.Default()
		logger.Println(err)
	}
}

-- name: ListDecks :many
SELECT * FROM decks 
ORDER BY name;

-- name: CreateDeck :one
INSERT INTO decks (
    name
) VALUES (
    $1
)
RETURNING *;
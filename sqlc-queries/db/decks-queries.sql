-- name: ListDecks :many
SELECT * FROM decks 
ORDER BY name;

-- name: GetDeck :one
SELECT * FROM decks 
WHERE name = $1;

-- name: CreateDeck :one
INSERT INTO decks (
    name
) VALUES (
    $1
)
RETURNING *;
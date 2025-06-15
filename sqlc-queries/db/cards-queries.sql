-- name: ListCards :many
SELECT * FROM cards
WHERE deck_id = $1
ORDER BY question;

-- name: CreateCard :one
INSERT INTO cards (
    deck_id, question, answer
) VALUES (
    $1, $2, $3
)
RETURNING *;
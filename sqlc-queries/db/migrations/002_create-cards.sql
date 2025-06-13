-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS cards (
  id   SERIAL PRIMARY KEY,
  question text    NOT NULL UNIQUE,
  answer text    NOT NULL,
  deck_id INTEGER NOT NULL CONSTRAINT fk_deck REFERENCES decks(id)
);

---- create above / drop below ----
DROP TABLE IF EXISTS cards;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

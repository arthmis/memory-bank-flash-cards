CREATE TABLE IF NOT EXISTS decks (
  id   SERIAL PRIMARY KEY,
  name text    NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS decks;

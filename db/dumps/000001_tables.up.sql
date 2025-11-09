

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS btree_gin;
CREATE EXTENSION IF NOT EXISTS plpgsql;

CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    parent_id INT REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    search_vector tsvector
);

CREATE OR REPLACE FUNCTION comments_tsvector_trigger() RETURNS trigger AS $$
BEGIN
  NEW.search_vector := to_tsvector('simple', unaccent(coalesce(NEW.content, '')));
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER comments_tsvector_update
BEFORE INSERT OR UPDATE
ON comments
FOR EACH ROW
EXECUTE FUNCTION comments_tsvector_trigger();

CREATE INDEX IF NOT EXISTS idx_comments_parent ON comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_comments_search ON comments USING GIN (search_vector);



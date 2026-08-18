ALTER TABLE movies ADD COLUMN
IF NOT EXISTS position INTEGER;
UPDATE movies SET position = id WHERE position IS NULL;

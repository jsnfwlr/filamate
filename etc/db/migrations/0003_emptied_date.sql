ALTER TABLE spools ADD COLUMN emptied_at TIMESTAMP WITH TIME ZONE NULL;
UPDATE spools SET emptied_at = NOW() WHERE empty = true;
ALTER TABLE spools DROP COLUMN empty;

---- create above / drop below ----

ALTER TABLE spools ADD COLUMN empty BOOLEAN NOT NULL DEFAULT false;
UPDATE spools SET empty = true WHERE emptied_at IS NOT NULL;
ALTER TABLE spools DROP COLUMN emptied_at;

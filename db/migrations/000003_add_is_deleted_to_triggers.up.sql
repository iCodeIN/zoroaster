BEGIN;

ALTER TABLE triggers ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT false;

COMMIT;
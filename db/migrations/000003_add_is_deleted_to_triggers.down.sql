BEGIN;

ALTER TABLE triggers DROP COLUMN is_deleted;

COMMIT;
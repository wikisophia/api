-- Delete the stuff created by create.sql
DROP INDEX IF EXISTS accounts_email_idx;
DROP TABLE IF EXISTS accounts;

DROP FUNCTION IF EXISTS update_last_modified;

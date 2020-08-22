/**
 * This file deletes all the database structures which were created
 * in create.sql.
 *
 * This file is run once for integration tests, immediately before
 * create.sql is run. All statements here must be guarded by IF EXISTS
 * clauses so that they can be run on an empty database without errors.
 *
 * It should be kept in sync with create.sql and empty.sql.
 */

DROP INDEX IF EXISTS argument_premises_premise_idx;
DROP INDEX IF EXISTS argument_premises_argument_version_idx;
DROP TABLE IF EXISTS argument_premises;

DROP INDEX IF EXISTS argument_versions_conclusion_idx;
DROP INDEX IF EXISTS argument_versions_argument_version_idx;
DROP INDEX IF EXISTS argument_versions_argument_idx;
DROP TABLE IF EXISTS argument_versions;

DROP TRIGGER IF EXISTS update_arguments_last_modified ON arguments;
DROP TABLE IF EXISTS arguments;

DROP INDEX IF EXISTS claims_claim_search_idx;
DROP INDEX IF EXISTS claims_claim_equals_idx;
DROP TABLE IF EXISTS claims;

DROP INDEX IF EXISTS accounts_email_idx;
DROP TABLE IF EXISTS accounts;

DROP FUNCTION IF EXISTS update_last_modified;

/**
 * This file deletes all the database structures which were created
 * in create.sql.
 *
 * This file is run once for integration tests, immediately before
 * create.sql is run. All statements here must be guarded by IF EXISTS
 * clauses so that they can be run on an empty database without errors.
 */

DROP TABLE IF EXISTS argument_premises;
DROP TABLE IF EXISTS argument_versions;
DROP TRIGGER IF EXISTS update_arguments_last_modified ON arguments;
DROP TABLE IF EXISTS arguments;
DROP TABLE IF EXISTS claims;
DROP FUNCTION IF EXISTS update_last_modified;

/**
 * This file wipes all the data from the DB without actually
 * destroying the structure.
 *
 * It's designed to be run at the start of every integration test case
 * which uses the database, to clear out state from the previous tests.
 *
 * It should be kept in sync with create.sql and destroy.sql.
 */

DELETE FROM argument_premises;
DELETE FROM argument_versions;
DELETE FROM arguments;
DELETE FROM claims;
DELETE FROM accounts;

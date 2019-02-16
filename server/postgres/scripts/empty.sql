/**
 * This file wipes all the data from the DB without actually
 * destroying the structure.
 *
 * It's designed to be run at the start of every integration test
 * which hits the database.
 */

DELETE FROM argument_premises;
DELETE FROM argument_versions;
DELETE FROM arguments;
DELETE FROM claims;

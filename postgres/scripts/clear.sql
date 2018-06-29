/**
 * This file initializes the database with all the needed tables.
 *
 * This is run during integration tests to clean the database between tests,
 * so it must be safe to run more than once.
 *
 * It may also be run to bootstrap the database on a new deploy.
 */

 /*
  * These commands are run to bootstrap the database. They aren't part of this script
  *
  * CREATE USER {config.postgres.user} WITH PASSWORD {config.postgres.password};
  * CREATE DATABASE {config.postgres.dbname} WITH OWNER {config.postgres.user}  ;
  */

DROP TABLE IF EXISTS premises;
DROP TABLE IF EXISTS arguments;
DROP TABLE IF EXISTS claims;

CREATE TABLE IF NOT EXISTS claims (
  id bigserial PRIMARY KEY,
  claim text UNIQUE NOT NULL CONSTRAINT claim_not_empty CHECK (claim != ''),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE claims IS 'Claims are used as both premises and conclusions of arguments.';
COMMENT ON COLUMN claims.claim IS 'This is the claim itself. For example, "Socrates is mortal".';
COMMENT ON COLUMN claims.created_at IS 'The time when this claim was added.';

-- Make the table insert-only
REVOKE ALL ON TABLE claims FROM PUBLIC;
GRANT INSERT ON TABLE claims TO app_wikisophia;


CREATE TABLE IF NOT EXISTS arguments (
  id bigserial PRIMARY KEY,
  conclusion_id bigint NOT NULL REFERENCES claims(id),
  latest_version smallint NOT NULL,
  live_version smallint NOT NULL CONSTRAINT live_less_max CHECK (live_version <= latest_version),
  deleted boolean NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE arguments IS 'There may be several arguments for the same conclusion. This table stores them all.';
COMMENT ON COLUMN arguments.conclusion_id IS 'The ID of the claim used as this argument''s conclusion.';
COMMENT ON COLUMN arguments.latest_version IS 'This tracks the number of revisions which have been made to this argument.';
COMMENT ON COLUMN arguments.live_version IS 'This tracks the version of the argument which is live on the site. In normal cases, this equals latest_version. It exists to support quick reverts to old versions.';
COMMENT ON COLUMN arguments.deleted IS 'If an argument was deleted, this is set to true. This allows deleted arguments to be restored easily.';
COMMENT ON COLUMN claims.created_at IS 'The time when this argument was first added.';

-- Make the table insert- and alter-able.
REVOKE ALL ON TABLE arguments FROM PUBLIC;
GRANT INSERT, UPDATE ON TABLE arguments TO app_wikisophia;


CREATE TABLE IF NOT EXISTS premises (
  id bigserial PRIMARY KEY,
  argument_id bigint NOT NULL REFERENCES arguments(id),
  argument_version smallint NOT NULL,
  claim_id bigint NOT NULL REFERENCES claims(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE premises IS 'This determines which claims are used as premises in a given argument.';
COMMENT ON COLUMN premises.argument_id IS 'The ID of the argument which this premise belongs to.';
COMMENT ON COLUMN premises.argument_version IS 'The version of the argument which this premise belongs to. This should be between 1 and arguments.latest_version, inclusive.';
COMMENT ON COLUMN premises.claim_id IS 'The ID of the claim used as this premise.';
COMMENT ON COLUMN claims.created_at IS 'The time when this premise was added.';

-- Make the table insert-only.
REVOKE ALL ON TABLE premises FROM PUBLIC;
GRANT INSERT ON TABLE premises TO app_wikisophia;

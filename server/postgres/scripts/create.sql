/**
 * This script initializes the database with all the tables
 * the app expects.
 *
 * This should be run by the default postgres user. It will create a user
 * with the same values as the app config.
 *
 * This gets run once at the start of the integration tests which use the database.
 * It's also intended for initializing a database for a new dev environment.
 *
 * Changes here must be kept in sync with destroy.sql and empty.sql.
 *
 * Note: this script expects the following to exist already:
 *   A user named "app_wikisophia" to exist already.
 *   A database named {config.postgres.dbname} WITH OWNER app_wikisophia
 */

CREATE FUNCTION update_last_modified()
RETURNS TRIGGER AS $$
BEGIN
   NEW.last_modified = NOW(); 
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE IF NOT EXISTS claims (
  id bigserial PRIMARY KEY,
  claim text UNIQUE NOT NULL CONSTRAINT claim_not_empty CHECK (claim != ''),
  created_on TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE claims IS 'Claims are used as both premises and conclusions of arguments.';
COMMENT ON COLUMN claims.claim IS 'This is the claim itself. For example, "Socrates is mortal".';
COMMENT ON COLUMN claims.created_on IS 'The time when this claim was added.';
REVOKE ALL ON TABLE claims FROM PUBLIC;
GRANT INSERT ON TABLE claims TO app_wikisophia;

CREATE TABLE IF NOT EXISTS arguments (
  id bigserial PRIMARY KEY,
  live_version integer NOT NULL DEFAULT 1,
  deleted boolean NOT NULL DEFAULT FALSE,
  created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_modified TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER update_arguments_last_modified BEFORE UPDATE ON arguments FOR EACH ROW EXECUTE PROCEDURE update_last_modified();
COMMENT ON TABLE arguments IS 'This stores all the arguments on wikisophia.';
COMMENT ON COLUMN arguments.live_version IS 'The version of the argument which is live. This refers to an argument_version in the argument_versions table.';
COMMENT ON COLUMN arguments.deleted IS 'True if this argument has been deleted, and false otherwise.';
COMMENT ON COLUMN arguments.created_on IS 'Timestamp of when the first version of this argument was created.';
COMMENT ON COLUMN arguments.last_modified IS 'Timestamp of when this argument was last modified/deleted/restored.';
REVOKE ALL ON TABLE arguments FROM PUBLIC;
GRANT INSERT, UPDATE ON TABLE arguments TO app_wikisophia;

CREATE TABLE argument_versions (
  id bigserial PRIMARY KEY,
  argument_id bigint NOT NULL REFERENCES arguments(id),
  argument_version integer NOT NULL,
  conclusion_id bigint NOT NULL REFERENCES claims(id),
  created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(argument_id, argument_version)
);
REVOKE ALL ON TABLE argument_versions FROM PUBLIC;
GRANT INSERT ON TABLE argument_versions TO app_wikisophia;
COMMENT ON TABLE argument_versions IS 'This tracks edits to an argument over time.';
COMMENT ON COLUMN argument_versions.argument_id IS 'The argument''s ID.';
COMMENT ON COLUMN argument_versions.argument_version IS 'The argument''s version.';
COMMENT ON COLUMN argument_versions.conclusion_id IS 'The argument''s conclusion.';
COMMENT ON COLUMN argument_versions.created_on IS 'Timestamp of when this version was created.';

CREATE TABLE argument_premises (
  id bigserial PRIMARY KEY,
  argument_version_id bigint NOT NULL REFERENCES argument_versions(id),
  premise_id bigint NOT NULL REFERENCES claims(id),
  created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(argument_version_id, premise_id)
);
REVOKE ALL ON TABLE argument_premises FROM PUBLIC;
GRANT INSERT ON TABLE argument_premises TO app_wikisophia;
COMMENT ON TABLE argument_premises IS 'This stores the premises used in each argument_version.';
COMMENT ON COLUMN argument_premises.argument_version_id IS 'The argument/version which includes this premise.';
COMMENT ON COLUMN argument_premises.premise_id IS 'The premise used in this version of the argument.';
COMMENT ON COLUMN argument_premises.created_on IS 'Timestamp of when this premise/version was created.';

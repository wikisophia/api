-- Create the stuff inside the arguments database.
-- Keep this in sync with the empty.sql and destroy.sql files.
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
CREATE INDEX claims_claim_equals_idx ON claims (claim);
CREATE INDEX claims_claim_search_idx ON claims USING gin(to_tsvector('english', claim));
REVOKE ALL ON TABLE claims FROM PUBLIC;
GRANT SELECT, INSERT ON TABLE claims TO :argumentsUser;

CREATE TABLE IF NOT EXISTS arguments (
  id bigserial PRIMARY KEY,
  deleted_on TIMESTAMPTZ DEFAULT NULL,
  created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_modified TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER update_arguments_last_modified BEFORE UPDATE ON arguments FOR EACH ROW EXECUTE PROCEDURE update_last_modified();
COMMENT ON TABLE arguments IS 'This stores all the arguments on wikisophia.';
COMMENT ON COLUMN arguments.deleted_on IS 'The timestamp when this argument was deleted. If null, it''s still live.';
COMMENT ON COLUMN arguments.created_on IS 'Timestamp of when the first version of this argument was created.';
COMMENT ON COLUMN arguments.last_modified IS 'Timestamp of when this argument was last deleted/restored. This does not update when the argument is edited to a new version.';
REVOKE ALL ON TABLE arguments FROM PUBLIC;
GRANT SELECT, INSERT, UPDATE ON TABLE arguments TO :argumentsUser;

CREATE TABLE argument_versions (
  id bigserial PRIMARY KEY,
  argument_id bigint NOT NULL REFERENCES arguments(id),
  argument_version integer NOT NULL,
  conclusion_id bigint NOT NULL REFERENCES claims(id),
  created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(argument_id, argument_version)
);
COMMENT ON TABLE argument_versions IS 'This tracks edits to an argument over time.';
COMMENT ON COLUMN argument_versions.argument_id IS 'The argument''s ID.';
COMMENT ON COLUMN argument_versions.argument_version IS 'The argument''s version.';
COMMENT ON COLUMN argument_versions.conclusion_id IS 'The argument''s conclusion.';
COMMENT ON COLUMN argument_versions.created_on IS 'Timestamp of when this version was created.';
CREATE INDEX argument_versions_argument_idx ON argument_versions (argument_id);
CREATE INDEX argument_versions_argument_version_idx ON argument_versions (argument_version);
CREATE INDEX argument_versions_conclusion_idx ON argument_versions (conclusion_id);
REVOKE ALL ON TABLE argument_versions FROM PUBLIC;
GRANT SELECT, INSERT ON TABLE argument_versions TO :argumentsUser;

CREATE TABLE argument_premises (
  id bigserial PRIMARY KEY,
  argument_version_id bigint NOT NULL REFERENCES argument_versions(id),
  premise_id bigint NOT NULL REFERENCES claims(id),
  created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(argument_version_id, premise_id)
);
COMMENT ON TABLE argument_premises IS 'This stores the premises used in each argument_version.';
COMMENT ON COLUMN argument_premises.argument_version_id IS 'The argument/version which includes this premise.';
COMMENT ON COLUMN argument_premises.premise_id IS 'The premise used in this version of the argument.';
COMMENT ON COLUMN argument_premises.created_on IS 'Timestamp of when this premise/version was created.';
CREATE INDEX argument_premises_argument_version_idx ON argument_premises (argument_version_id);
CREATE INDEX argument_premises_premise_idx ON argument_premises (premise_id);
REVOKE ALL ON TABLE argument_premises FROM PUBLIC;
GRANT SELECT, INSERT ON TABLE argument_premises TO :argumentsUser;

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO :argumentsUser;

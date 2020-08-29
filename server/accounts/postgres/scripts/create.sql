-- Create the stuff inside the accounts database.
-- Keep this in sync with the empty.sql and destroy.sql files.
CREATE FUNCTION update_last_modified()
RETURNS TRIGGER AS $$
BEGIN
   NEW.last_modified = NOW();
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE IF NOT EXISTS accounts (
  id bigserial PRIMARY KEY,
  email varchar(100) NOT NULL,
  reset_token varchar(100),
  reset_token_expiry TIMESTAMPTZ,
  password_hash varchar(5000),
  created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(email),
  CONSTRAINT tokens_must_expire CHECK (reset_token IS NULL OR reset_token_expiry IS NOT NULL)
);
CREATE TRIGGER update_accounts_last_modified BEFORE UPDATE ON accounts FOR EACH ROW EXECUTE PROCEDURE update_last_modified();
COMMENT ON TABLE accounts IS 'The accounts that exist on the site.';
COMMENT ON COLUMN accounts.email IS 'The email associated with this account. Each account has a unique email.';
COMMENT ON COLUMN accounts.reset_token IS 'The token generated by the app to reset this account''s password. This will be null if the user hasn''t requested a reset recently.';
COMMENT ON COLUMN accounts.reset_token_expiry IS 'The timestamp when the password reset_token expires.';
COMMENT ON COLUMN accounts.password_hash IS 'The hashed password value. This may be null if the email hasn''t been verified yet';
COMMENT ON COLUMN accounts.created_on IS 'Timestamp of when this account was created.';
COMMENT ON COLUMN accounts.last_modified IS 'Timestamp of when this row was last modified.';
CREATE INDEX accounts_email_idx ON accounts (email);
REVOKE ALL ON TABLE accounts FROM PUBLIC;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE accounts TO :accountsUser;

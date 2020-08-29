-- Make the databases/users for the accounts database.

CREATE USER app_wikisophia_accounts_test WITH PASSWORD 'app_wikisophia_accounts_test_password';
CREATE DATABASE wikisophia_accounts_test WITH OWNER app_wikisophia_accounts_test;

CREATE USER :accountsUser WITH PASSWORD :accountsPass;
CREATE DATABASE wikisophia_accounts WITH OWNER :accountsUser;

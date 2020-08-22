-- Run this file once, as the postgres user, to initialize the database.

CREATE USER app_wikisophia_test WITH PASSWORD 'app_wikisophia_test_password';
CREATE DATABASE wikisophia_test WITH OWNER app_wikisophia_test;

CREATE USER app_wikisophia WITH PASSWORD 'app_wikisophia_password';
CREATE DATABASE wikisophia WITH OWNER app_wikisophia;

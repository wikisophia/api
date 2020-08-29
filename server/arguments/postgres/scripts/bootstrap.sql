-- Make the databases/users for the arguments database.

CREATE USER app_wikisophia_arguments_test WITH PASSWORD 'app_wikisophia_arguments_test_password';
CREATE DATABASE wikisophia_arguments_test WITH OWNER app_wikisophia_arguments_test;

CREATE USER :argumentsUser WITH PASSWORD :argumentsPass;
CREATE DATABASE wikisophia_arguments WITH OWNER :argumentsUser;

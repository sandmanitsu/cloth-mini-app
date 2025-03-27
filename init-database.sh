#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "admin" --dbname "storage" <<-EOSQL
    CREATE DATABASE users;
    GRANT ALL PRIVILEGES ON DATABASE storage TO admin;
    GRANT ALL PRIVILEGES ON DATABASE users TO admin;
EOSQL
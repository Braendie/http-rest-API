CREATE TABLE users (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  id_telegram INT UNIQUE,
  email VARCHAR UNIQUE,
  encrypted_password VARCHAR
);
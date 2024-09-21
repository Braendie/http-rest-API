CREATE TABLE users (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  id_telegram INT UNIQUE,
  email VARCHAR NOT NULL UNIQUE,
  encrypted_password VARCHAR NOT NULL,
  height INT NOT NULL,
  age INT NOT NULL,
  weight INT NOT NULL,
  gender VARCHAR(6) NOT NULL,
  phone_number VARCHAR UNIQUE,
);
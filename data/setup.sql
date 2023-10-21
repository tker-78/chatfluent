drop table sessions;
drop table posts;
drop table threads;
drop table users;

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  uuid VARCHAR(64) NOT NULL UNIQUE,
  name VARCHAR(255),
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
  id SERIAL PRIMARY KEY,
  uuid VARCHAR(64) NOT NULL UNIQUE,
  email VARCHAR(225),
  user_id INTEGER REFERENCES users(id),
  created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS threads (
  id SERIAL PRIMARY KEY,
  uuid VARCHAR(64) NOT NULL UNIQUE,
  topic VARCHAR(255),
  user_id INTEGER REFERENCES users(id),
  created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
  id SERIAL PRIMARY KEY,
  uuid VARCHAR(64) NOT NULL UNIQUE,
  body TEXT,
  user_id INTEGER REFERENCES users(id),
  thread_id INTEGER REFERENCES threads(id),
  created_at TIMESTAMP NOT NULL
);
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS intents (
  intent_label VARCHAR(32) PRIMARY KEY NOT NULL,
  -- embedding vector(4096) NOT NULL
  embedding vector(1536) NOT NULL
);

CREATE TABLE IF NOT EXISTS requests (
  event_id UUID NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  source VARCHAR(255) NOT NULL,
  client_id VARCHAR(255) NOT NULL,
  latitude DOUBLE PRECISION NOT NULL,
  longitude DOUBLE PRECISION NOT NULL,
  message_text VARCHAR(4095) NOT NULL
);

CREATE TABLE IF NOT EXISTS responses (
  event_id UUID NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  message_text VARCHAR(4095) NOT NULL,
  media_url VARCHAR(511) NOT NULL,
  media_type VARCHAR(255) NOT NULL,
  action VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS pushes (
  event_id UUID NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  client_id VARCHAR(255) NOT NULL,
  message_text VARCHAR(4095) NOT NULL,
  media_url VARCHAR(511) NOT NULL,
  media_type VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS scheduled_events (
  event_id UUID NOT NULL,
  event_type VARCHAR(255) NOT NULL,
  scheduled TIMESTAMP NOT NULL,
  created TIMESTAMP NOT NULL,
  info JSON NOT NULL
);


CREATE INDEX IF NOT EXISTS idx_scheduled_events_event_type ON scheduled_events(event_type,scheduled);

CREATE TABLE IF NOT EXISTS clients_registry (
  source VARCHAR(255) NOT NULL,
  user_id UUID,
  client_id VARCHAR(255) NOT NULL,
  info JSON NOT NULL,
  latitude FLOAT8 NOT NULL,
  longitude FLOAT8 NOT NULL,
  created TIMESTAMP NOT NULL,
  updated TIMESTAMP NOT NULL,
  PRIMARY KEY (source, client_id)
);

CREATE TABLE IF NOT EXISTS users_registry (
  user_id UUID NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  created TIMESTAMP NOT NULL,
  updated TIMESTAMP NOT NULL
);

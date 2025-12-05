CREATE TABLE events(
  event_id VARCHAR PRIMARY KEY,
  event_type VARCHAR NOT NULL,
  processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  completed_at TIMESTAMPTZ
);

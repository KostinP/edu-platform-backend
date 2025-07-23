CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  telegram_id TEXT UNIQUE,
  first_name TEXT,
  last_name TEXT,
  username TEXT,
  photo_url TEXT,
  email TEXT,
  subscribe_to_newsletter BOOLEAN DEFAULT FALSE,
  role TEXT NOT NULL DEFAULT 'unspecified',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
  visitor_id UUID NULL
);

CREATE TABLE visitor_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    visitor_id UUID NOT NULL, -- это UUID из куки, 
    event_type TEXT NOT NULL, -- строка, описывающая событие, 
    event_data JSONB, -- дополнительные данные в JSON (например, url, user agent и т.п.).
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_visitor_events_visitor_id ON visitor_events(visitor_id);
CREATE INDEX idx_visitor_events_created_at ON visitor_events(created_at);

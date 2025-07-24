-- Расширения
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Таблица пользователей
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
  inactivity_expiry_days INTEGER DEFAULT 180, -- пользователь может указать, через сколько дней неактивности токен недействителен
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP,
  visitor_id UUID NULL
);

-- События визиторов (до авторизации)
CREATE TABLE visitor_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visitor_id UUID NOT NULL,
  event_type TEXT NOT NULL,
  event_data JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_visitor_events_visitor_id ON visitor_events(visitor_id);
CREATE INDEX idx_visitor_events_created_at ON visitor_events(created_at);

-- Сессии пользователей
CREATE TABLE user_sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token TEXT NOT NULL,
  user_agent TEXT,
  ip_address TEXT,
  country TEXT,
  city TEXT,
  created_at TIMESTAMP DEFAULT now(),
  last_active_at TIMESTAMP DEFAULT now(),
  expires_at TIMESTAMP,
  revoked_at TIMESTAMP, -- для удаления сессий
  is_current BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_last_active_at ON user_sessions(last_active_at);
CREATE INDEX idx_user_sessions_token ON user_sessions(token);

CREATE TABLE IF NOT EXISTS user_inactivity_timeout (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    timeout_seconds BIGINT NOT NULL
);
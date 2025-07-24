DROP TABLE IF EXISTS user_inactivity_timeout;

DROP INDEX IF EXISTS idx_user_sessions_user_id;
DROP INDEX IF EXISTS idx_user_sessions_last_active_at;
DROP INDEX IF EXISTS idx_user_sessions_token;
DROP TABLE IF EXISTS user_sessions;

DROP INDEX IF EXISTS idx_visitor_events_visitor_id;
DROP INDEX IF EXISTS idx_visitor_events_created_at;
DROP TABLE IF EXISTS visitor_events;

DROP TABLE IF EXISTS users;

-- 00001_new_users_table_down.sql

DROP TRIGGER IF EXISTS update_users_update_at ON users;

DROP FUNCTION IF EXISTS update_update_at_column();

DROP TABLE IF EXISTS users;

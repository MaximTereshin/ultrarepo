-- Удаление триггеров
DROP TRIGGER IF EXISTS update_wallets_updated_at ON wallets;
DROP TRIGGER IF EXISTS update_games_updated_at ON games;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаление индексов
DROP INDEX IF EXISTS idx_wallets_user_id;
DROP INDEX IF EXISTS idx_transactions_wallet_id;
DROP INDEX IF EXISTS idx_transactions_type;
DROP INDEX IF EXISTS idx_transactions_status;
DROP INDEX IF EXISTS idx_game_sessions_user_id;
DROP INDEX IF EXISTS idx_game_sessions_game_id;
DROP INDEX IF EXISTS idx_game_sessions_status;
//fsdfasdfs
-- Удаление таблиц
DROP TABLE IF EXISTS game_sessions;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS wallets;

//бабабебе


//пепепепупупу

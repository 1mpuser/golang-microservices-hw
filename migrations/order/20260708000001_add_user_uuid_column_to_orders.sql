-- +goose Up
-- Шаг 1/3: добавить nullable-колонку user_uuid (без NOT NULL — иначе упадёт на существующих строках).
-- TODO(неделя 5, часть 2): ALTER TABLE orders ADD COLUMN user_uuid UUID;
ALTER TABLE orders ADD COLUMN user_uuid UUID

-- +goose Down
-- TODO(неделя 5): ALTER TABLE orders DROP COLUMN IF EXISTS user_uuid;
ALTER TABLE orders drop COLUMN user_uuid
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages (
	id BIGSERIAL PRIMARY KEY,
	recept_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
	pacient_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
	message_text TEXT NOT NULL,
	is_read BOOLEAN  DEFAULT FALSE,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_chat_thread ON messages(recept_id, pacient_id, created_at);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_chat_thread;
DROP TABLE messages;
-- +goose StatementEnd

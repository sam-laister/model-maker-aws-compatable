-- +goose Up
-- +goose StatementBegin

-- Drop foreign key constraint from chat_messages to chats
ALTER TABLE chats
DROP CONSTRAINT IF EXISTS fk_chat_task;

ALTER TABLE chat_messages
DROP CONSTRAINT IF EXISTS fk_message_chat;

-- Drop the chats table
DROP TABLE IF EXISTS chats;

-- Add task_id to chat_messages
ALTER TABLE chat_messages
ADD COLUMN task_id INTEGER NOT NULL;

-- Add foreign key constraint from chat_messages to tasks
ALTER TABLE chat_messages
ADD CONSTRAINT fk_chat_messages_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop the foreign key constraint and column

-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin

-- Optional: Remove the user_id column if it exists
ALTER TABLE chats
DROP COLUMN IF EXISTS user_id;

-- Add the task_id column and create a foreign key constraint
ALTER TABLE chats
ADD COLUMN task_id INTEGER NOT NULL;

ALTER TABLE chats
ADD CONSTRAINT fk_chat_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd

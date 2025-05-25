-- +goose Up
-- +goose StatementBegin
-- Create the TaskStatus enum
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'taskstatus') THEN
        CREATE TYPE TaskStatus AS ENUM ('SUCCESS', 'INPROGRESS', 'FAILED', 'INITIAL');
    END IF;
END $$;

-- Create the tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    title TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    status TaskStatus NOT NULL DEFAULT 'INITIAL',
    user_id INT NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the app_files table

-- Drop the tasks table
DROP TABLE IF EXISTS tasks;

-- Drop the TaskStatus enum
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'taskstatus') THEN
        DROP TYPE TaskStatus;
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin

-- Create the collections table
CREATE TABLE IF NOT EXISTS collections (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT fk_collection_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create the collection_tasks join table
CREATE TABLE IF NOT EXISTS collection_tasks (
    collection_id INT NOT NULL,
    task_id INT NOT NULL,
    PRIMARY KEY (collection_id, task_id),
    CONSTRAINT fk_collection FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE,
    CONSTRAINT fk_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop the collection_tasks join table
DROP TABLE IF EXISTS collection_tasks;

-- Drop the collections table
DROP TABLE IF EXISTS collections;

-- +goose StatementEnd

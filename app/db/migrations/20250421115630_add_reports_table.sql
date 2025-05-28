-- +goose Up
-- +goose StatementBegin

-- Create the ReportType enum
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'reporttype') THEN
        CREATE TYPE ReportType AS ENUM ('BUG', 'FEEDBACK');
    END IF;
END $$;

-- Create the reports table
CREATE TABLE IF NOT EXISTS reports (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    report_type ReportType NOT NULL,
    rating INTEGER,
    user_id INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_report_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop the reports table
DROP TABLE IF EXISTS reports;

-- Drop the ReportType enum
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'reporttype') THEN
        DROP TYPE ReportType;
    END IF;
END $$;

-- +goose StatementEnd

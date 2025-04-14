-- +goose Up
-- +goose StatementBegin
CREATE TABLE StatusTask (
    id SERIAL PRIMARY KEY,
    status_task TEXT NOT NULL
);

INSERT INTO StatusTask(status_task) 
VALUES ('CREATED'), ('PROCESSING'), ('FAILED'), ('COMPLETED'); 

CREATE TABLE TasksLog (
    id BIGSERIAL PRIMARY KEY,
    audit_log JSON NOT NULL,
    status_task INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    FOREIGN KEY(status_task) REFERENCES StatusTask(id) 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS TasksLog;
DROP TABLE IF EXISTS StatusTask;
-- +goose StatementEnd

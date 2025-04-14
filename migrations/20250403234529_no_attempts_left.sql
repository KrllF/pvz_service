-- +goose Up
-- +goose StatementBegin
INSERT INTO StatusTask(status_task)
VALUES('NO_ATTEMPTS_LEFT');

ALTER TABLE TasksLog
ADD attempts INT NOT NULL DEFAULT 0;

UPDATE TasksLog
SET attempts = 1
WHERE status_task != 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM TasksLog
WHERE status_task = (
    SELECT id FROM StatusTask WHERE status_task = 'NO_ATTEMPTS_LEFT'
);

DELETE FROM StatusTask
WHERE status_task = 'NO_ATTEMPTS_LEFT';

ALTER TABLE TasksLog
DROP COLUMN attempts;
-- +goose StatementEnd
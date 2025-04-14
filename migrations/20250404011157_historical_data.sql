-- +goose Up
-- +goose StatementBegin
INSERT INTO tasksLog(audit_log, status_task, created_at)
SELECT 
    json_build_object(
        'method', method,
        'request', request,
        'code', code
    ),
    (SELECT id FROM statustask WHERE status_task='CREATED'),
    created_at
    FROM handlerlog
    WHERE created_at < '2025-03-29 20:41:51.400';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM tasksLog
WHERE created_at < '2025-03-29 20:41:51.400';
-- +goose StatementEnd

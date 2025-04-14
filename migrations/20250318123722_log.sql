-- +goose Up
-- +goose StatementBegin
CREATE TABLE HandlerLog (
    method TEXT NOT NULL,
    request TEXT NOT NULL,
    code INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE StatusLog (
    order_id BIGINT NOT NULL,
    new_status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS HandlerLog;
DROP TABLE IF EXISTS StatusLog;
-- +goose StatementEnd

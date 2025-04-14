-- +goose Up
-- +goose StatementBegin
CREATE TABLE Users (
    user_id BIGINT PRIMARY KEY
);

CREATE TABLE Pack (
    pack_id BIGSERIAL PRIMARY KEY,
    pack_type TEXT NOT NULL,
    is_extra BOOLEAN NOT NULL
);

CREATE TABLE StatusOrder (
    status_id BIGSERIAL PRIMARY KEY,
    status_type TEXT NOT NULL
);

CREATE TABLE Orders (
    order_id BIGINT PRIMARY KEY ,
    user_id BIGINT NOT NULL,
    order_weight BIGINT NOT NULL,
    order_price BIGINT NOT NULL,
    pack_id BIGINT NOT NULL,
    extra_pack_id BIGINT,
    status_id BIGINT NOT NULL,
    shelf_life TIMESTAMP NOT NULL,
    two_days_of_life TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY(user_id) REFERENCES Users(user_id),
    FOREIGN KEY(pack_id) REFERENCES Pack(pack_id),
    FOREIGN KEY(extra_pack_id) REFERENCES Pack(pack_id),
    FOREIGN KEY(status_id) REFERENCES StatusOrder(status_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Orders;
DROP TABLE IF EXISTS StatusOrder;
DROP TABLE IF EXISTS Pack;
DROP TABLE IF EXISTS Users;
-- +goose StatementEnd

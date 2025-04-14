-- +goose Up
-- +goose StatementBegin
INSERT INTO Pack (pack_type, is_extra)
VALUES('film', TRUE), ('box', FALSE), ('bag', FALSE);

INSERT INTO StatusOrder (status_type)
VALUES('delivered'), ('accepted'), ('returned');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM Orders
WHERE pack_id IN (
    SELECT pack_id FROM Pack WHERE pack_type IN ('film', 'box', 'bag')
)  OR status_id IN (
    SELECT status_id FROM StatusOrder WHERE status_type IN ('delivered', 'accepted', 'returned')
);

DELETE FROM Pack
WHERE pack_type IN ('film', 'box', 'bag');

DELETE FROM StatusOrder
WHERE status_type IN ('delivered', 'accepted', 'returned');
-- +goose StatementEnd

-- +goose Up
INSERT INTO wallets (id, balance)
VALUES
    ('1c63a43f-aacd-47b0-bc3b-535e69c6ed4c', 100000),
    ('c3a8cb84-03f2-4fb9-982a-9ee2cfb50b9f', 50000),
    ('c01841a5-368b-4901-951e-ec067ab4c4fa', 0);

-- +goose Down
DELETE FROM wallets
WHERE id IN (
    '1c63a43f-aacd-47b0-bc3b-535e69c6ed4c',
    'c3a8cb84-03f2-4fb9-982a-9ee2cfb50b9f',
    'c01841a5-368b-4901-951e-ec067ab4c4fa'
);
use `isucari`;

ALTER TABLE items ADD INDEX ix_seller_status_created_at(seller_id, status, created_at);

INSERT INTO `public_items` (item) SELECT id from `items` ORDER BY created_at ASC, id ASC;

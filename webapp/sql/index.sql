use `isucari`;

ALTER TABLE items ADD INDEX ix_seller_status_created_at(seller_id, status, created_at);
ALTER TABLE items ADD INDEX (created_at);

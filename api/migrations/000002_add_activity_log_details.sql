-- +goose Up
ALTER TABLE activity_logs ADD COLUMN entity_type VARCHAR(255);
ALTER TABLE activity_logs ADD COLUMN entity_id UUID;

-- +goose Down
ALTER TABLE activity_logs DROP COLUMN entity_id;
ALTER TABLE activity_logs DROP COLUMN entity_type;

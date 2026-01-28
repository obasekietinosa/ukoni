-- +goose Up
ALTER TABLE invitations ADD COLUMN token VARCHAR(255);
UPDATE invitations SET token = uuid_generate_v4()::text WHERE token IS NULL;
ALTER TABLE invitations ALTER COLUMN token SET NOT NULL;
CREATE UNIQUE INDEX invitations_token_idx ON invitations(token);

-- +goose Down
DROP INDEX IF EXISTS invitations_token_idx;
ALTER TABLE invitations DROP COLUMN IF EXISTS token;

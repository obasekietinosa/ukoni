-- +goose Up
CREATE TABLE password_reset_tokens (
    token VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);

-- +goose Down
DROP TABLE IF EXISTS password_reset_tokens;

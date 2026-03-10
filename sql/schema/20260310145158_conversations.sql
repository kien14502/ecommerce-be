-- +goose Up
-- +goose StatementBegin
CREATE TABLE conversations (
    id CHAR(36) PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_conversation_member_user 
ON conversation_members(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `conversations`;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE reactions (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36),
    post_id CHAR(36),
    reaction_type VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, post_id),

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);
CREATE INDEX idx_reactions_post 
ON reactions(post_id);

CREATE INDEX idx_reactions_user 
ON reactions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `reactions`;
-- +goose StatementEnd

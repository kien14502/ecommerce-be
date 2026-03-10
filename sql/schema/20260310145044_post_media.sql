-- +goose Up
-- +goose StatementBegin
CREATE TABLE post_media (
    id CHAR(36) PRIMARY KEY,
    post_id CHAR(36) NOT NULL,
    media_url TEXT,
    media_type VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `post_media`;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE otp_verifications (
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    otp_hash TEXT NOT NULL,
    purpose ENUM('register','login','reset_password','verify_email'),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_otp_email 
ON otp_verifications(email);

CREATE INDEX idx_otp_expire 
ON otp_verifications(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `otp_verifications`;
-- +goose StatementEnd

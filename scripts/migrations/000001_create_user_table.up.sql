CREATE table users
(
    id                   BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    login                VARCHAR(255)    NOT NULL UNIQUE,
    inn                  BIGINT UNSIGNED NOT NULL UNIQUE,
    name                 VARCHAR(100)    NOT NULL,
    password             VARCHAR(255)    NOT NULL,
    account_provider_id  VARCHAR(255) UNIQUE,
    customer_provider_id VARCHAR(255) UNIQUE,
    created_at           TIMESTAMP       NOT NULL,
    updated_at           TIMESTAMP       NOT NULL,
    deleted_at           TIMESTAMP       NULL,
    is_banned            boolean DEFAULT FALSE
);
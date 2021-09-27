CREATE table users
(
    id         INT          NOT NULL PRIMARY KEY AUTO_INCREMENT,
    login      VARCHAR(255) NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL,
    deleted_at TIMESTAMP    NULL,
    is_banned  boolean DEFAULT TRUE
);
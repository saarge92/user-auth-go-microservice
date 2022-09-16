CREATE TABLE currencies
(
    id          INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    code        VARCHAR(4) UNIQUE,
    description VARCHAR(255) NULL,
    created_at  TIMESTAMP    NOT NULL,
    updated_at  TIMESTAMP    NOT NULL
);

INSERT INTO currencies (id, code, description, created_at, updated_at)
VALUES (1, 'RUB', 'Russian rouble', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP()),
       (2, 'USD', 'American dollar', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP()),
       (3, 'EUR', 'European currency', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP())

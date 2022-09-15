CREATE table users
(
    id                   BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    login                VARCHAR(255)    NOT NULL UNIQUE,
    inn                  VARCHAR(255)    NOT NULL UNIQUE,
    name                 VARCHAR(100)    NOT NULL,
    password             VARCHAR(255)    NOT NULL,
    account_provider_id  VARCHAR(255) UNIQUE,
    customer_provider_id VARCHAR(255) UNIQUE,
    created_at           INT UNSIGNED    NOT NULL,
    updated_at           INT UNSIGNED    NOT NULL,
    deleted_at           INT UNSIGNED    NULL,
    is_banned            boolean DEFAULT FALSE,
    country_id           INT UNSIGNED    NULL,
    CONSTRAINT `player_country_id_fk` FOREIGN KEY (country_id) REFERENCES countries (id)
        ON DELETE SET NULL ON UPDATE CASCADE
);

INSERT INTO users
SET id                   = 1,
    login                = 'user@foo.com',
    inn                  = '7721546864',
    name                 = 'ivan',
    password             = '$2a$12$1H/oEHp4xCv8KMQFVVeNR.NQEKCyzs9jjGSEWyvoi8MNyINmFbtgK', #qwerty123`
    account_provider_id  = 'account-uuid',
    customer_provider_id = 'customer-uuid',
    created_at           = unix_timestamp(now()),
    updated_at           = unix_timestamp(now()),
    country_id           = 1;
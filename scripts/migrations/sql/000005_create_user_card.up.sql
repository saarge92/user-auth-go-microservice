CREATE TABLE cards
(
    id                   INT UNSIGNED PRIMARY KEY NOT NULL AUTO_INCREMENT,
    external_id          CHAR(36)                 NOT NULL UNIQUE,
    number               VARCHAR(255)             NOT NULL,
    expire_month         SMALLINT UNSIGNED        NOT NULL,
    expire_year          INT UNSIGNED             NOT NULL,
    user_id              BIGINT UNSIGNED          NOT NULL,
    external_provider_id VARCHAR(255)             NOT NULL,
    is_default           BOOLEAN                       DEFAULT FALSE,
    created_at           INT UNSIGNED             NOT NULL,
    updated_at           INT UNSIGNED             NOT NULL,
    deleted_at           INT UNSIGNED             NULL DEFAULT NULL,
    CONSTRAINT `card_user_id_fk` FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE NO ACTION
        ON UPDATE CASCADE,
    CONSTRAINT `unique_card_user` UNIQUE (number, user_id, external_provider_id)
);
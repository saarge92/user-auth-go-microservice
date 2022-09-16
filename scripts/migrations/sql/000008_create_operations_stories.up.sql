CREATE TABLE operations_stories
(
    id                   BIGINT UNSIGNED     NOT NULL PRIMARY KEY AUTO_INCREMENT,
    external_id          CHAR(36) UNIQUE     NOT NULL,
    user_id              BIGINT UNSIGNED     NOT NULL,
    card_id              INT UNSIGNED        NOT NULL,
    amount               DECIMAL(36, 18)     NOT NULL,
    balance_before       DECIMAL(36, 18)     NOT NULL,
    balance_after        DECIMAL(36, 18)     NOT NULL,
    external_provider_id VARCHAR(255) UNIQUE NOT NULL,
    operation_type_id    INT UNSIGNED        NOT NULL,
    created_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),

    CONSTRAINT `fk_operation_stories_user_id` FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE CASCADE,
    CONSTRAINT `fk_operation_stories_card_id` FOREIGN KEY (card_id) REFERENCES cards (id)
        ON DELETE NO ACTION ON UPDATE CASCADE,
    CONSTRAINT `fk_operation_stories_operation_type_id` FOREIGN KEY (operation_type_id)
        REFERENCES operation_types (id)
);
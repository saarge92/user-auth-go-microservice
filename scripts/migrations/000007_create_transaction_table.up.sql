CREATE TABLE transactions
(
    id                   INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    external_id          CHAR(36)        NOT NULL UNIQUE,
    external_provider_id VARCHAR(255)    NOT NULL UNIQUE,
    wallet_id_from       BIGINT UNSIGNED NOT NULL,
    wallet_id_to         BIGINT UNSIGNED NOT NULL,
    exchange_rate        DECIMAL         NOT NULL,
    created_at           TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    CONSTRAINT `fk_transaction_wallet_id_from` FOREIGN KEY (wallet_id_from) REFERENCES wallets (id)
        ON DELETE NO ACTION ON UPDATE CASCADE,
    CONSTRAINT `fk_transaction_wallet_id_to` FOREIGN KEY (wallet_id_to) REFERENCES wallets (id)
);
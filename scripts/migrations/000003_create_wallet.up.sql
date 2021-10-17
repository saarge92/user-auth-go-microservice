CREATE TABLE wallets
(
    id          BIGINT unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id     BIGINT UNSIGNED NOT NULL,
    currency_id INT UNSIGNED    NOT NULL,
    balance     DECIMAL(36, 13) NOT NULL,
    is_default  BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP       NOT NULL,
    updated_at  TIMESTAMP       NOT NULL,
    CONSTRAINT `fk_wallet_to_user_id` FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE NO ACTION ON UPDATE CASCADE,
    CONSTRAINT `fk_wallet_to_currency_id` FOREIGN KEY (currency_id)
        REFERENCES currencies (id) ON DELETE NO ACTION ON UPDATE CASCADE
);
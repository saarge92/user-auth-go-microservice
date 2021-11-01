CREATE TABLE countries
(
    id            INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    code_2        VARCHAR(2)   NOT NULL,
    code_3        VARCHAR(3)   NOT NULL,
    currency_code VARCHAR(5)   NOT NULL,
    phone_code    INT UNSIGNED NOT NULL
);

ALTER TABLE users
    ADD COLUMN country_id INT UNSIGNED NULL;

ALTER TABLE users
    ADD CONSTRAINT `player_country_id_fk` FOREIGN KEY (country_id) REFERENCES countries (id)
        ON DELETE SET NULL ON UPDATE CASCADE;

INSERT INTO countries (code_2, code_3, currency_code, phone_code)
VALUES ('US', 'USA', 'USD', 1),
       ('RU', 'RUS', 'RUB', 7)

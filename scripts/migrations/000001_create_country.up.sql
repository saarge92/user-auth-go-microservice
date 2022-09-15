CREATE TABLE countries
(
    id            INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    code_2        VARCHAR(2)   NOT NULL,
    code_3        VARCHAR(3)   NOT NULL,
    currency_code VARCHAR(5)   NOT NULL,
    phone_code    INT UNSIGNED NOT NULL
);

INSERT INTO countries (id, code_2, code_3, currency_code, phone_code)
VALUES (1, 'US', 'USA', 'USD', 1),
       (2, 'RU', 'RUS', 'RUB', 7)

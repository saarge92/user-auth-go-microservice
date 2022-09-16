CREATE TABLE operation_types
(
    id   INT UNSIGNED PRIMARY KEY,
    name varchar(255) UNIQUE NOT NULL
);

INSERT into operation_types(id, name)
VALUES (1, 'Deposit'),
       (2, 'Refund');
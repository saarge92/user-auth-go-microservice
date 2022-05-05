CREATE TABLE roles
(
    id   INT UNSIGNED NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

INSERT INTO roles(id, name)
VALUES (1, 'User'), (2, 'Admin');

CREATE TABLE IF NOT EXISTS user_roles
(
    id      INT UNSIGNED    NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    role_id INT UNSIGNED NOT NULL,
    CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles (id),
    CONSTRAINT unique_user_role_ids UNIQUE (user_id, role_id)
);
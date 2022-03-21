CREATE TABLE user
(
    id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    city VARCHAR(255),
    country VARCHAR(255),
    date_of_birth DATE,
    email VARCHAR(255),
    name VARCHAR(100),
    postal_code VARCHAR(5),
    surname VARCHAR(100),
    test_alter_type INT,
    test_alter_rename INT
);

ALTER TABLE user ADD example VARCHAR(255);

ALTER TABLE user DROP COLUMN country;

ALTER TABLE user ALTER COLUMN test_alter_type TYPE VARCHAR(255);

ALTER TABLE user RENAME COLUMN test_alter_rename TO test_alter_rename_new;

CREATE TABLE book
(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100),
    user_id BIGINT REFERENCES user(id)
);

CREATE TABLE a
(
    id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT
);

CREATE TABLE b
(
    id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT
);

CREATE TABLE a_to_b
(
    a_id BIGINT REFERENCES a(id),
    b_id BIGINT REFERENCES b(id),
    example VARCHAR(255),
    PRIMARY KEY(a_id, b_id)
);
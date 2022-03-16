CREATE TABLE user
(
    id INT PRIMARY KEY NOT NULL,
    city VARCHAR(255),
    country VARCHAR(255),
    date_of_birth DATE,
    email VARCHAR(255),
    name VARCHAR(100),
    postal_code VARCHAR(5),
    surname VARCHAR(100)
);

CREATE TABLE book
(
    id INT PRIMARY KEY NOT NULL,
    name VARCHAR(100),
    user_id INT REFERENCES user(id)
);

CREATE TABLE a
(
    id INT PRIMARY KEY NOT NULL
);

CREATE TABLE b
(
    id INT PRIMARY KEY NOT NULL
);

CREATE TABLE a_to_b
(
    a_id INT REFERENCES a(id),
    b_id INT REFERENCES b(id),
    example VARCHAR(255),
    PRIMARY KEY(a_id, b_id)
);
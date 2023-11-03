CREATE TABLE my_table (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    age INTEGER
);

INSERT INTO my_table (name, age)
VALUES ('John Doe', 28),
       ('Jane Smith', 32),
       ('Bob Johnson', 45);
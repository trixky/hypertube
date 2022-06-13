CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR (30) NOT NULL,
    firstname VARCHAR (30) NOT NULL,
    lastname VARCHAR (30) NOT NULL,
    email VARCHAR (90) NOT NULL,
    password VARCHAR (65) NOT NULL
);
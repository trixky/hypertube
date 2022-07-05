CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR (30) NOT NULL,
    firstname VARCHAR (30) NOT NULL,
    lastname VARCHAR (30) NOT NULL,
    email VARCHAR (90) UNIQUE NOT NULL,
    id_42 INTEGER UNIQUE,
    password VARCHAR (65)
);
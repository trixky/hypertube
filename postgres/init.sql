\c hypertube;

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR (30) NOT NULL,
    firstname VARCHAR (30) NOT NULL,
    lastname VARCHAR (30) NOT NULL,
    email VARCHAR (90) UNIQUE NOT NULL,
    id_42 INTEGER UNIQUE,
    password VARCHAR (65)
);

INSERT INTO users (username, firstname, lastname, email, password) VALUES ('admin', 'mathis', 'bois', 'mathis@email.com', 'c54b18a947c806a48d7fd825ec6aea73');
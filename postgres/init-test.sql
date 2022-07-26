\ c hypertube;

-- unencrypted from server side password: 1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6
-- **************** TESTS META
-- user test 0 (0)
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (0, 'username_test_0', 'firstname_test_0', 'lastname_test_0', 'email.test_0@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');

-- **************** TESTS API-AUTH
-- user test 1 (1)
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (101, 'username_test_1', 'firstname_test_1', 'lastname_test_1', 'email.test_1@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');
-- user test 2 (2)
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (102, 'username_test_2', 'firstname_test_2', 'lastname_test_2', 'email.test_2@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');
-- user test 3 (3)
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (103, 'username_test_3', 'firstname_test_3', 'lastname_test_3', 'email.test_3@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');
-- user test 4 (4)
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (104, 'username_test_4', 'firstname_test_4', 'lastname_test_4', 'email.test_4@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');
-- user test 5 (5) (reserved for recover password apply)
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (105, 'username_test_5', 'firstname_test_5', 'lastname_test_5', 'email.test_5@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');

-- **************** TESTS API-USER
-- user test 1 (6)
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (201, 'username_test_6', 'firstname_test_6', 'lastname_test_6', 'email.test_6@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');
-- user test 2 (7)
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (202, 'username_test_7', 'firstname_test_7', 'lastname_test_7', 'email.test_7@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');
-- user test 3 (8)
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (203, 'username_test_8', 'firstname_test_8', 'lastname_test_8', 'email.test_8@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');
-- user test 4 (9)
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (204, 'username_test_9', 'firstname_test_9', 'lastname_test_9', 'email.test_9@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');

-- **************** TESTS API-MEDIA
INSERT INTO medias (id, imdb_id, tmdb_id, description, duration, thumbnail, background, year, rating)
VALUES (1, 'tt0000000', 0, 'Movie movie movie', 180, NULL, NULL, 1990, 6.9)
INsERT INTO media_names (id, media_id, name, lang)
VALUES (1, 1, 'Movie', '__')

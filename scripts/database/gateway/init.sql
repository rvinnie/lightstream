CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    path character varying(255) NOT NULL
);

INSERT INTO images (id, path) VALUES
(1, 'images/eg.jpg'),
(2, 'images/kapibara.png');
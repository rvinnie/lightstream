CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    path character varying(255) NOT NULL
);

INSERT INTO images (path) VALUES
('images/eg.jpg'),
('images/kapibara.png');
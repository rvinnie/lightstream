CREATE TABLE history (
    id SERIAL PRIMARY KEY,
    videoId INT NOT NULL,
    watched TIMESTAMP NOT NULL
);
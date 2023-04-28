-- Tables read by sqlc to generate models.go due to special syntax.
--
-- This file is NOT READ BY THE DATABASE.

CREATE TABLE terms (
    text TEXT NOT NULL,
    variants TEXT NOT NULL,
    part_of_speech TEXT NOT NULL,
    common_level TEXT NOT NULL,
    translations TEXT NOT NULL,
    popularity TEXT NOT NULL
);
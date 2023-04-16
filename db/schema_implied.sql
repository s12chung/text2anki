-- Tables read by sqlc to generate models.go due to special syntax.
--
-- This file is NOT READ BY THE DATABASE.

CREATE TABLE terms (
    text text NOT NULL,
    variants text NOT NULL,
    part_of_speech text NOT NULL,
    common_level text NOT NULL,
    translations text NOT NULL,
    popularity text NOT NULL
);
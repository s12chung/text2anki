CREATE TABLE IF NOT EXISTS terms (
    text TEXT NOT NULL,
    variants TEXT NOT NULL,
    part_of_speech TEXT NOT NULL,
    common_level INTEGER NOT NULL,
    translations TEXT NOT NULL,
    popularity INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS sources (
    id INTEGER PRIMARY KEY,
    tokenized_texts TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
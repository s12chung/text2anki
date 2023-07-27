CREATE TABLE terms (
    id INTEGER PRIMARY KEY,
    text TEXT NOT NULL,
    part_of_speech TEXT NOT NULL,

    variants TEXT NOT NULL,
    translations TEXT NOT NULL,
    common_level INTEGER NOT NULL,

    popularity INTEGER NOT NULL
);

CREATE TABLE sources (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    parts TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TRIGGER sources_updated_at
    BEFORE UPDATE ON sources FOR EACH ROW
BEGIN
    UPDATE sources SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

CREATE TABLE notes (
    id INTEGER PRIMARY KEY,

    text TEXT NOT NULL,
    part_of_speech TEXT NOT NULL,
    translation TEXT NOT NULL,
    explanation TEXT NOT NULL,

    common_level INTEGER NOT NULL,
    usage TEXT NOT NULL,
    usage_translation TEXT NOT NULL,
    dictionary_source TEXT NOT NULL,

    notes TEXT NOT NULL,

    downloaded BOOLEAN DEFAULT false NOT NULL
);
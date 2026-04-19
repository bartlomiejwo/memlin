CREATE TABLE languages (
    id SERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL
);

CREATE TABLE words (
    id SERIAL PRIMARY KEY,
    language_id INTEGER REFERENCES languages(id),
    word TEXT NOT NULL,
    base_form_id INTEGER REFERENCES words(id),
    pronunciation TEXT,
    category TEXT,
    level TEXT,
    popularity FLOAT
);

CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    word_id INTEGER REFERENCES words(id),
    translation_id INTEGER REFERENCES words(id),
    UNIQUE (word_id, translation_id)
);

CREATE TABLE grammatical_forms (
    id SERIAL PRIMARY KEY,
    word_id INTEGER REFERENCES words(id),
    base_word_id INTEGER REFERENCES words(id),
    grammar_type TEXT NOT NULL,
    grammar_notes TEXT
);

CREATE TABLE contexts (
    id SERIAL PRIMARY KEY,
    word_id INTEGER REFERENCES words(id),
    language_id INTEGER REFERENCES languages(id),
    context TEXT NOT NULL
);

CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    word_id INTEGER REFERENCES words(id),
    context_id INTEGER REFERENCES contexts(id),
    media_type TEXT NOT NULL,
    path TEXT NOT NULL,
    description TEXT,
    user_id INTEGER,
    CHECK (word_id IS NOT NULL OR context_id IS NOT NULL)
);
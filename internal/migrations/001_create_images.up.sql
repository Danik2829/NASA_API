CREATE TABLE metadata(
    date DATE NOT NULL UNIQUE,
    explanation TEXT NOT NULL,
    hdurl TEXT NOT NULL,
    media_type TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL
);
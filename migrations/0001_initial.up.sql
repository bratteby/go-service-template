CREATE TABLE example (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

GRANT SELECT, INSERT, UPDATE, DELETE ON example to example;
-- Список всех тегов
CREATE TABLE tags (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    author_id UUID NOT NULL,

    name TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE tag_assignments (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    author_id UUID NOT NULL,

    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    target_type TEXT NOT NULL,              -- 'lesson', 'course', 'module', etc.
    target_id UUID NOT NULL
);


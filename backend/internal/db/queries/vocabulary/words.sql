-- name: GetWord :one
SELECT w.id, w.word, l.name AS language, w.pronunciation, w.category, w.level, w.popularity
FROM words w
JOIN languages l ON w.language_id = l.id
WHERE w.id = $1;

-- name: ListWords :many
SELECT w.id, w.word, l.name AS language, w.pronunciation, w.category, w.level, w.popularity
FROM words w
JOIN languages l ON w.language_id = l.id
ORDER BY w.id
LIMIT $1 OFFSET $2;

-- name: CreateWord :one
INSERT INTO words (
    language_id, 
    word, 
    pronunciation, 
    category, 
    level, 
    popularity
) VALUES (
    (SELECT id FROM languages WHERE code = $1),
    $2, $3, $4, $5, $6
)
RETURNING id, word, pronunciation, category, level, popularity;

-- name: UpdateWord :one
UPDATE words
SET 
    word = $2,
    pronunciation = $3,
    category = $4,
    level = $5,
    popularity = $6
WHERE id = $1
RETURNING id, word, pronunciation, category, level, popularity;

-- name: DeleteWord :exec
DELETE FROM words
WHERE id = $1;
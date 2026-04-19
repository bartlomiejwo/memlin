-- name: GetLanguage :one
SELECT id, code, name
FROM languages
WHERE id = $1;

-- name: GetLanguageByCode :one
SELECT id, code, name
FROM languages
WHERE code = $1;

-- name: ListLanguages :many
SELECT id, code, name
FROM languages
ORDER BY name;

-- name: CreateLanguage :one
INSERT INTO languages (
    code, name
) VALUES (
    $1, $2
)
RETURNING id, code, name;
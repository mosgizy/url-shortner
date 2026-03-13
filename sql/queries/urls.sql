-- name: CreateUrl :one
INSERT INTO urls (id, short_code, long_url, clicks, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetByShortCode :one
SELECT id, short_code, long_url, created_at, updated_at, clicks 
FROM urls 
WHERE short_code = $1;

-- name: IncrementClicks :exec
UPDATE urls SET clicks = clicks + 1, updated_at = NOW()
WHERE short_code = 
-- name: CreateLink :one
INSERT INTO links (short_code, long_url) VALUES ($1, $2) RETURNING id, short_code, long_url, clicks;

-- name: GetLinkByCode :one
SELECT * FROM links WHERE short_code = $1 LIMIT 1;
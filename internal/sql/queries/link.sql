-- name: CreateLink :one
INSERT INTO links (short_code, long_url) VALUES ($1, $2) RETURNING id, short_code, long_url, clicks;

-- name: GetClicksByCode :one
SELECT clicks FROM links WHERE short_code = $1 LIMIT 1;

-- name: GetURLAndIncrementLinkClicks :one
UPDATE links SET clicks = clicks + 1 WHERE short_code = $1 RETURNING long_url;
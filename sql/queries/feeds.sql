-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id)
VALUES(
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name, feeds.url, users.name from feeds
INNER JOIN users
ON
feeds.user_id = users.id;

-- name: GetFeedByURL :one
SELECT * from feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = NOW(),
last_fetched_at = NOW()
where id = $1;

-- name: GetNextFeedToFetch :one
SELECT * from feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
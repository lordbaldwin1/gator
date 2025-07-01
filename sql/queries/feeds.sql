-- name: AddFeed :one
INSERT INTO feeds (id, name, url, created_at, updated_at, user_id)
VALUES(
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT name, url, user_id FROM feeds;

-- name: GetFeedsWithUsername :many
SELECT feeds.name, feeds.url, users.name
FROM feeds
INNER JOIN users 
ON feeds.user_id = users.id;

-- name: GetFeedByUserID :one
SELECT *
FROM feeds
WHERE user_id = $1
LIMIT 1;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE url = $1
LIMIT 1;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    values($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT
  inserted_feed_follow.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN users on users.id = inserted_feed_follow.user_id;

-- name: GetFeedFollowsForUser :many
SELECT users.name, feeds.name
FROM users
INNER JOIN feed_follows
ON users.id = feed_follows.user_id
INNER JOIN feeds
ON feeds.id = feed_follows.feed_id
WHERE users.id = $1;

-- name: DeleteFeedFollow :one
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1
AND feed_follows.feed_id
IN (SELECT feed_id FROM feeds WHERE feeds.url = $2)
RETURNING *;
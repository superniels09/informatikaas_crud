-- name: CreateComment :execresult
-- desc: Create a new comment on a game
INSERT INTO comments (
    content,
    user_id,
    game_id
) VALUES (
    ?, ?, ?
);

-- name: GetCommentByID :one
-- desc: Get a comment by its ID
SELECT * FROM comments WHERE id = ?;

-- name: GetCommentsByGameID :many
-- desc: Get all comments for a specific game (newest first)
SELECT * FROM comments
WHERE game_id = ?
ORDER BY created_at DESC;

-- name: GetCommentsByUserID :many
-- desc: Get all comments by a specific user (newest first)
SELECT * FROM comments
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: ListComments :many
-- desc: List all comments with pagination (newest first)
SELECT * FROM comments
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateComment :execresult
-- desc: Update a comment's content
UPDATE comments
SET content = ?
WHERE id = ?;

-- name: DeleteComment :exec
-- desc: Delete a comment by its ID
DELETE FROM comments WHERE id = ?;

-- name: GetCommentWithUserAndGame :one
-- desc: Get a comment with its associated user and game information
SELECT
    c.*,
    u.username AS user_username,
    u.display_name AS user_display_name,
    g.title AS game_title
FROM comments c
JOIN users u ON c.user_id = u.id
JOIN games g ON c.game_id = g.id
WHERE c.id = ?;

-- name: GetGameCommentsCount :one
-- desc: Get the count of comments for a specific game
SELECT COUNT(*) FROM comments WHERE game_id = ?;

-- name: GetUserCommentsCount :one
-- desc: Get the count of comments by a specific user
SELECT COUNT(*) FROM comments WHERE user_id = ?;

-- name: DeleteGameComments :exec
-- desc: Delete all comments for a specific game (cascade delete)
DELETE FROM comments WHERE game_id = ?;

-- name: DeleteUserComments :exec
-- desc: Delete all comments by a specific user (cascade delete)
DELETE FROM comments WHERE user_id = ?;

-- name: ListGameCommentsWithPagination :many
-- desc: List comments for a specific game with pagination
SELECT * FROM comments
WHERE game_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

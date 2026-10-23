-- name: CreateGame :execresult
INSERT INTO games (
    title,
    description,
    script,
    user_id
) VALUES (
    ?, ?, ?, ?
);

-- name: GetGameByID :one
-- desc: Get a game by its ID
SELECT * FROM games WHERE id = ?;

-- name: GetGamesByUserID :many
-- desc: Get all games by user ID
SELECT * FROM games WHERE user_id = ? ORDER BY created_at DESC;

-- name: ListGames :many
-- desc: List all games (with pagination)
SELECT * FROM games
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateGame :execresult
-- desc: Update a game's information
UPDATE games
SET
    title = COALESCE(?, title),
    description = COALESCE(?, description),
    script = COALESCE(?, script),
    user_id = COALESCE(?, user_id)
WHERE id = ?;

-- name: UpdateGameScript :execresult
-- desc: Update only a game's script
UPDATE games
SET script = ?
WHERE id = ?;

-- name: DeleteGame :exec
-- desc: Delete a game by its ID
DELETE FROM games WHERE id = ?;

-- name: GetGameWithUser :one
-- desc: Get a game with its associated user information
SELECT
    g.*,
    u.username AS user_username,
    u.display_name AS user_display_name
FROM games g
JOIN users u ON g.user_id = u.id
WHERE g.id = ?;

-- name: GetUserGamesCount :one
-- desc: Get the count of games for a specific user
SELECT COUNT(*) FROM games WHERE user_id = ?;

-- name: DeleteUserGames :exec
-- desc: Delete all games associated with a user (cascade delete)
DELETE FROM games WHERE user_id = ?;

-- name: AddFavorite :exec
-- desc: Add a favorite relationship between user and game
INSERT INTO favorites (user_id, game_id)
VALUES (?, ?);

-- name: RemoveFavorite :exec
-- desc: Remove a favorite relationship between user and game
DELETE FROM favorites
WHERE user_id = ? AND game_id = ?;

-- name: GetUserFavoritedGames :many
-- desc: Get all games favorited by a specific user
SELECT
    g.*,
    f.created_at AS favorited_at
FROM favorites f
JOIN games g ON f.game_id = g.id
WHERE f.user_id = ?
ORDER BY f.created_at DESC;

-- name: IsGameFavorited :one
-- desc: Check if a user has favorited a specific game
SELECT EXISTS(
    SELECT 1 FROM favorites
    WHERE user_id = ? AND game_id = ?
) AS is_favorited;

-- name: GetGameFavoritedByUsers :many
-- desc: Get all users who favorited a specific game
SELECT
    u.*,
    f.created_at AS favorited_at
FROM favorites f
JOIN users u ON f.user_id = u.id
WHERE f.game_id = ?
ORDER BY f.created_at DESC;

-- name: GetGameFavoriteCount :one
-- desc: Get the total number of favorites for a specific game
SELECT COUNT(*) FROM favorites WHERE game_id = ?;

-- name: GetUserFavoriteCount :one
-- desc: Get the total number of favorites by a specific user
SELECT COUNT(*) FROM favorites WHERE user_id = ?;

-- name: GetUserFavoritedGamesCount :one
-- desc: Get the count of games favorited by a specific user
SELECT COUNT(*) FROM favorites WHERE user_id = ?;

-- name: GetGamesFavoriteCountWithPagination :many
-- desc: Get games with their favorite counts (paginated)
SELECT
    g.*,
    (SELECT COUNT(*) FROM favorites WHERE game_id = g.id) AS favorite_count
FROM games g
ORDER BY favorite_count DESC, g.created_at DESC
LIMIT ? OFFSET ?;

-- name: GetUserFavoritesWithGameInfo :many
-- desc: Get all games favorited by a user with game details
SELECT
    g.*,
    f.created_at AS favorited_at,
    u.username AS user_username
FROM favorites f
JOIN games g ON f.game_id = g.id
JOIN users u ON f.user_id = u.id
WHERE f.user_id = ?
ORDER BY f.created_at DESC;

-- name: CreateUser :execresult
-- desc: Create a new user
INSERT INTO users (
    username,
    passhash,
    email,
    display_name,
    bio
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: GetUserByID :one
-- desc: Get a user by their ID
SELECT * FROM users WHERE id = ?;

-- name: GetUserByUsername :one
-- desc: Get a user by their username
SELECT * FROM users WHERE username = ?;

-- name: GetUserByEmail :one
-- desc: Get a user by their email
SELECT * FROM users WHERE email = ?;

-- name: ListUsers :many
SELECT * FROM users
LIMIT ? OFFSET ?;

-- name: UpdateUser :execresult
-- desc: Update a user's information
UPDATE users
SET
    username = COALESCE(?, username),
    passhash = COALESCE(?, passhash),
    email = COALESCE(?, email),
    display_name = COALESCE(?, display_name),
    bio = COALESCE(?, bio)
WHERE id = ?;

-- name: UpdateUserPassword :execresult
-- desc: Update only a user's password
UPDATE users
SET passhash = ?
WHERE id = ?;

-- name: DeleteUser :exec
-- desc: Delete a user by their ID
DELETE FROM users WHERE id = ?;

-- name: GetUserCount :one
-- desc: Get the total number of users
SELECT COUNT(*) FROM users;

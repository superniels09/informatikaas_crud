-- name: MessagesGet :one
SELECT * FROM messages
WHERE id = ? LIMIT 1;

-- name: MessagesList :many
SELECT * FROM messages;

-- name: MessagesCreate :execresult
INSERT INTO messages (
  content
) VALUES (
  ?
);

-- name: MessagesDelete :exec
DELETE FROM messages
WHERE id = ?;

-- name: MessagesDeleteAll :exec
DELETE FROM messages;

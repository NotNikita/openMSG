-- name: CreateUser :one
INSERT INTO users (nickname, public_key, avatar)
VALUES ($1, $2, $3)
RETURNING id, nickname, public_key, avatar, created_at;

-- name: GetUserByID :one
SELECT id, nickname, public_key, avatar, created_at
FROM users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT id, nickname, public_key, avatar, created_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetOrCreateConversation :one
INSERT INTO conversations (user_a_id, user_b_id)
VALUES ($1, $2)
ON CONFLICT (user_a_id, user_b_id) DO UPDATE SET user_a_id = EXCLUDED.user_a_id
RETURNING id, user_a_id, user_b_id, created_at;

-- name: GetConversationsByUserID :many
SELECT id, user_a_id, user_b_id, created_at
FROM conversations
WHERE user_a_id = $1 OR user_b_id = $1
ORDER BY created_at DESC;

-- name: GetConversationByID :one
SELECT id, user_a_id, user_b_id, created_at
FROM conversations
WHERE id = $1;

-- name: CreateMessage :one
INSERT INTO messages (conversation_id, sender_id, ciphertext, nonce)
VALUES ($1, $2, $3, $4)
RETURNING id, conversation_id, sender_id, ciphertext, nonce, created_at;

-- name: GetMessagesByConversationID :many
SELECT id, conversation_id, sender_id, ciphertext, nonce, created_at
FROM messages
WHERE conversation_id = $1
ORDER BY created_at ASC;

-- name: GetPublicMessages :many
SELECT
    u_sender.nickname AS sender_nickname,
    u_recip.nickname  AS recipient_nickname,
    m.ciphertext,
    m.created_at
FROM messages m
JOIN conversations c ON m.conversation_id = c.id
JOIN users u_sender ON m.sender_id = u_sender.id
JOIN users u_recip ON (
    (c.user_a_id = m.sender_id AND u_recip.id = c.user_b_id) OR
    (c.user_b_id = m.sender_id AND u_recip.id = c.user_a_id)
)
WHERE (sqlc.narg('before')::timestamptz IS NULL OR m.created_at < sqlc.narg('before'))
ORDER BY m.created_at DESC
LIMIT sqlc.arg('limit');

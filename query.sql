--
-- AccessToken Table
--

-- name: GetAccessToken :one
SELECT * FROM tokens
WHERE id = $1 LIMIT 1;

-- name: FindAccessTokenByHash :one
SELECT * FROM tokens
WHERE hash = $1 LIMIT 1;

-- name: CreateAccessToken :one
INSERT INTO tokens
(user_id,
 hash,
 expires_at,
 last_used_at,
 created_utc,
 updated_utc)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: DeleteAccessToken :exec
DELETE FROM tokens
WHERE id = $1;


--
-- EmailLog Table
--

-- name: CreateEmailLog :exec
INSERT INTO email_logs
(user_id, message_type, created_at)
VALUES ($1, $2, NOW());

-- name: CountRecentEmails :one
SELECT count(*)
FROM email_logs
WHERE user_id = $1
    AND message_type = $2
    AND created_at >= NOW() - INTERVAL '31 days';


--
-- users Table
--

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: FindUserByEmployeeID :one
SELECT *
FROM users
WHERE employee_id = $1
LIMIT 1;

-- name: FindUserByUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: FindUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users
(employee_id, first_name, last_name, display_name, username, email, "active", "locked",
 last_login_at, updated_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6, TRUE, FALSE, NOW(), NOW(), NOW()) RETURNING *;

-- name: ListActiveUnlockedUsers :many
SELECT * FROM users WHERE active and NOT locked;

-- name: FindUsersToPurge :many
SELECT *
FROM users
WHERE NOT active AND updated_at < $1;

-- name: UpdateUserLastLoggedIn :exec
UPDATE users
SET last_login_at   = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET employee_id          = $2,
    first_name           = $3,
    last_name            = $4,
    display_name         = $5,
    username             = $6,
    email                = $7,
    active               = $8,
    locked               = $9,
    updated_at           = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

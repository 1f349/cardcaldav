-- name: GetPasswordHash :one
SELECT password
FROM mailbox
WHERE username = ?
  AND active > 0;

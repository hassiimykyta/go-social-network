-- name: CreateMedia :one
INSERT INTO media (
  owner_id, kind, storage_key, mime_type, size_bytes, width, height, duration_ms
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING media.*;

-- name: AttachMediaToPost :exec
INSERT INTO post_media (post_id, media_id, position)
VALUES ($1, $2, $3)
ON CONFLICT (post_id, media_id)
DO UPDATE SET position = EXCLUDED.position;
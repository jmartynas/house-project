-- name: SelectSutartys :many
SELECT * FROM sutartis;

-- name: SelectBreziniai :many
SELECT * FROM brezinys;

-- name: SelectLeidimai :many
SELECT * FROM leidimas;

-- name: SelectSutartis :one
SELECT * FROM sutartis
WHERE id = $1 LIMIT 1;

-- name: SelectBrezinys :one
SELECT * FROM brezinys
WHERE id = $1 LIMIT 1;

-- name: SelectLeidimas :one
SELECT * FROM leidimas
WHERE id = $1 LIMIT 1;

-- name: InsertSutartis :exec
INSERT INTO sutartis (
sutartis,
kaina
) VALUES (
$1, $2
);

-- name: InsertBrezinys :exec
INSERT INTO brezinys (
brezinys,
fk_sutartis_id
) VALUES (
$1, $2
);

-- name: InsertLeidimas :exec
INSERT INTO leidimas (
leidimas,
fk_brezinys_id
) VALUES (
$1, $2
);

-- name: UpdateSutartis :exec
UPDATE sutartis SET
sutartis = $2,
kaina = $3
WHERE id = $1;

-- name: UpdateBrezinys :exec
UPDATE brezinys SET
brezinys = $2,
fk_sutartis_id = $3
WHERE id = $1;

-- name: UpdateLeidimas :exec
UPDATE leidimas SET
leidimas = $2,
fk_brezinys_id = $3
WHERE id = $1;

-- name: DeleteSutartis :exec
DELETE FROM sutartis
WHERE id = $1;

-- name: DeleteBrezinys :exec
DELETE FROM brezinys
WHERE id = $1;

-- name: DeleteLeidimas :exec
DELETE FROM leidimas
WHERE id = $1;

-- name: DeleteLeidimai :exec
DELETE FROM leidimas
WHERE fk_brezinys_id = $1;

-- name: SelectDeleteLeidimai :many
SELECT * FROM leidimas
WHERE fk_brezinys_id = $1;

-- name: SelectDeleteBreziniai :many
SELECT * FROM brezinys
WHERE fk_sutartis_id = $1;

-- name: CheckSutartys :one
SELECT * FROM sutartis_id_seq;

-- name: CheckBreziniai :one
SELECT * FROM brezinys_id_seq;

-- name: CheckLeidimai :one
SELECT * FROM leidimas_id_seq;

-- name: GetBehaviorsOfClient :many
SELECT * FROM behaviors 
WHERE cst_dim_id = ?;
--
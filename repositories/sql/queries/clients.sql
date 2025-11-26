-- name: GetClients :many
SELECT * FROM clients;
--

-- name: GetClientByCstDimId :one
SELECT * FROM clients 
WHERE cst_dim_id = ?;
--
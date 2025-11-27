-- name: GetTransactions :many
SELECT  c.first_name, c.last_name, 
    t.transdatetime, t.amount, t.direction, t.target 
    FROM transactions AS t
LEFT JOIN clients AS c
ON t.cst_dim_id = c.cst_dim_id;
--

-- name: GetTransactionsOfClient :many
SELECT  transdatetime, amount, direction, target 
    FROM transactions 
WHERE cst_dim_id = ?;
--
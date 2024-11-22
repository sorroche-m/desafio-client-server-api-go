-- name: CadastrarCotacaoDolar :exec
INSERT INTO cotacao_dolar (valor, data) VALUES (?, ?);

-- name: BuscarCotacaoDolar :many
SELECT * FROM cotacao_dolar ORDER BY data DESC LIMIT 1;

-- name: BuscarCotacaoDolarPorData :one
SELECT * FROM cotacao_dolar WHERE data = ?;
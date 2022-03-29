package sql

import (
	"database/sql"
	"net/url"
)

// SQLStr ...
type SQLStr struct {
	url *url.URL
	db  *sql.DB
}

type OPS []Op

//Op ...
type Op struct {
	Ref     string  `json:"referencia,omitempty"`
	Ean     *string `json:"ean13,omitempty"`
	Nome    string  `json:"nome,omitempty"`
	Cor     string  `json:"cor,omitempty"`
	Tamanho string  `json:"tamanho,omitempty"`
	Uni     string  `json:"unidade,omitempty"`
	Quanti  string  `json:"quantidade,omitempty"`
	Ex1     string  `json:"extra1,omitempty"`
	Ex2     string  `json:"extra2,omitempty"`
	Ex20    string  `json:"extra20,omitempty"`
	Grupo   string  `json:"grupo,omitempty"`
}

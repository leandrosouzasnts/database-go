package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {

	stringConexao := "root:password@/devbook?charset=utf8&parseTime=True&Local"

	db, erro := sql.Open("mysql", stringConexao)

	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil
}

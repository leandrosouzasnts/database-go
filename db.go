package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // É uma declaração implicita no qual fica armazenado de forma externa
)

func main() {

	//URL := "usuario:senha/nome_banco"

	stringConexao := "root:password@/devbook?charset=utf8&parseTime=True&Local"

	db, erro := sql.Open("mysql", stringConexao)

	if erro != nil {
		log.Fatal(erro)
	}

	defer db.Close()

	if erro = db.Ping(); erro != nil {
		log.Fatal(erro)
	}

	fmt.Println("Connection success")

	rows, erro := db.Query("SELECT * FROM users")

	if erro != nil {
		log.Fatal(erro)
	}

	defer rows.Close()

	fmt.Println(rows)
}

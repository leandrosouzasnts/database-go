package server

import (
	"database-go/db"
	"database-go/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CreatedUser(w http.ResponseWriter, r *http.Request) {

	requestBody, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		w.Write([]byte("Falha ao ler corpo da requisição"))
		return
	}

	var user model.User

	if erro := json.Unmarshal(requestBody, &user); erro != nil {
		w.Write([]byte("Falha ao converter JSON para Struct User"))
		return
	}

	connection, erro := db.ConnectDB()
	if erro != nil {
		w.Write([]byte("Falha ao conectar no Banco de Dados: " + erro.Error()))
		return
	}
	defer connection.Close()

	stmt, erro := connection.Prepare("INSERT INTO users (name) VALUES (?)")
	if erro != nil {
		w.Write([]byte(erro.Error()))
		return
	}
	defer stmt.Close()

	// _, erro = stmt.Exec(user.Name)
	// if erro != nil {
	// 	w.Write([]byte(erro.Error()))
	// 	return
	// }

	transaction, erro := connection.Begin()
	if erro != nil {
		w.Write([]byte(erro.Error()))
	}

	_, erro = transaction.Stmt(stmt).Exec(user.Name)

	if erro != nil {
		transaction.Rollback()
		w.Write([]byte("Falha ao Gravar Usuário no Banco de Dados"))
	}

	erro = transaction.Commit()
	if erro != nil {
		w.Write([]byte(erro.Error()))
	}

	fmt.Printf("Usuário: %s cadastrado no Banco de Dados", user.Name)

	var userResponse model.User
	err := connection.QueryRow("SELECT id, name FROM users WHERE id = LAST_INSERT_ID()").Scan(&userResponse.ID, &userResponse.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao recuperar usuário inserido no Banco de Dados: " + err.Error()))
		return
	}

	responseBody, erro := json.Marshal(userResponse)
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao criar resposta JSON"))
		return
	}

	// Defina o tipo de conteúdo da resposta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Escreva a resposta JSON na resposta HTTP
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)

}

func GetUser(w http.ResponseWriter, r *http.Request) {

}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	connection, error := db.ConnectDB()
	if error != nil {
		w.Write([]byte("Falha ao conectar no Banco de Dados"))
		return
	}

	rows, error := connection.Query("SELECT id, name FROM users")
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Falha ao carregar usuários do Banco de Dados"))
		return
	}

	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User
		error := rows.Scan(&user.ID, &user.Name)
		if error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Falha ao atribuir usuário do Banco de Dados"))
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Escreva a resposta JSON na resposta HTTP
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

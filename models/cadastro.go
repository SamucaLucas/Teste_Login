package models

import (
	"app_login/config" // Lembre-se de usar o nome do módulo que está no seu go.mod
	"log"
)

// CriarUsuario insere um novo registro no banco de dados
func CriarUsuario(nome, email, senha string) error {
	db := config.ConectarBD()
	defer db.Close() // Fecha a conexão quando a função terminar

	// Prepara a query SQL para evitar SQL Injection
	query := "INSERT INTO usuarios (nome, email, senha) VALUES ($1, $2, $3)"
	inserirDados, err := db.Prepare(query)
	if err != nil {
		log.Println("Erro ao preparar a query de insert:", err)
		return err
	}
	defer inserirDados.Close()

	// Executa a query com os valores recebidos
	_, err = inserirDados.Exec(nome, email, senha)
	if err != nil {
		log.Println("Erro ao executar o insert (E-mail já existe?):", err)
		return err
	}

	return nil
}

package main

import (
	"app_login/controllers" // Ajuste para o nome do seu módulo go.mod
	"fmt"
	"net/http"
)

func main() {
	// Definindo as Rotas
	http.HandleFunc("/", controllers.TelaLogin)
	http.HandleFunc("/cadastro", controllers.TelaCadastro)
	http.HandleFunc("/dashboard", controllers.Dashboard)

	http.HandleFunc("/login", controllers.ProcessarLogin)
	http.HandleFunc("/logout", controllers.Logout)

	http.HandleFunc("/editar", controllers.Editar)
	http.HandleFunc("/update", controllers.Atualizar)
	http.HandleFunc("/insert", controllers.ProcessarCadastro)

	http.HandleFunc("/deletar", controllers.Deletar)

	fmt.Println("Servidor rodando na porta 8080...")
	fmt.Println("Acesse: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

package models

import (
	"app_login/config"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// ValidarLogin verifica se o e-mail existe e se a senha está correta
func ValidarLogin(campo1, campo2 string) (Usuario, error) {
	db := config.ConectarBD()
	defer db.Close()

	usuario := Usuario{}
	var senhaHash string

	// TENTATIVA 1: O comportamento normal (campo1 = email, campo2 = senha)
	err := db.QueryRow("SELECT id, nome, email, senha FROM usuarios WHERE email=$1", campo1).Scan(&usuario.Id, &usuario.Nome, &usuario.Email, &senhaHash)
	if err == nil {
		// Encontrou o e-mail no campo1. Agora verifica se o campo2 é a senha correta
		if errBcrypt := bcrypt.CompareHashAndPassword([]byte(senhaHash), []byte(campo2)); errBcrypt == nil {
			return usuario, nil // Login normal com sucesso
		}
	}

	// TENTATIVA 2 (O BUG): Se falhou, tenta inverter! (campo2 = email, campo1 = senha)
	err = db.QueryRow("SELECT id, nome, email, senha FROM usuarios WHERE email=$1", campo2).Scan(&usuario.Id, &usuario.Nome, &usuario.Email, &senhaHash)
	if err == nil {
		// Encontrou o e-mail no campo2! Agora verifica se o campo1 é a senha correta
		if errBcrypt := bcrypt.CompareHashAndPassword([]byte(senhaHash), []byte(campo1)); errBcrypt == nil {
			return usuario, nil // LOGIN COM SUCESSO MESMO INVERTIDO! (Defeito de lógica)
		}
	}

	// Se as duas combinações falharem, aí sim retorna erro
	return usuario, errors.New("credenciais invalidas")
}

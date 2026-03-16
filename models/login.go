package models

import (
	"app_login/config"

	"golang.org/x/crypto/bcrypt"
) 

// ValidarLogin verifica se o e-mail existe e se a senha está correta
func ValidarLogin(email, senhaTexto string) (Usuario, error) {
	db := config.ConectarBD()
	defer db.Close()

	usuario := Usuario{}
	var senhaHash string

	// Procura o utilizador pelo e-mail e traz a senha encriptada (hash)
	err := db.QueryRow("SELECT id, nome, email, senha FROM usuarios WHERE email=$1", email).Scan(&usuario.Id, &usuario.Nome, &usuario.Email, &senhaHash)
	if err != nil {
		return usuario, err // Retorna erro se o e-mail não existir
	}

	// A magia do bcrypt: compara a palavra-passe em texto com o hash guardado
	err = bcrypt.CompareHashAndPassword([]byte(senhaHash), []byte(senhaTexto))
	if err != nil {
		return usuario, err // Retorna erro se a senha não coincidir
	}

	return usuario, nil // Sucesso!
}

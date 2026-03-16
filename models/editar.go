package models

import "app_login/config"

// BuscarUsuario retorna um único usuário pelo ID para preencher o formulário
func BuscarUsuario(id string) Usuario {
	db := config.ConectarBD()
	defer db.Close()

	// Busca apenas o usuário que tenha o ID passado na URL
	resultado, err := db.Query("SELECT id, nome, email FROM usuarios WHERE id=$1", id)
	if err != nil {
		panic(err.Error())
	}
	defer resultado.Close()

	usuario := Usuario{}

	for resultado.Next() {
		var id int
		var nome, email string

		err = resultado.Scan(&id, &nome, &email)
		if err != nil {
			panic(err.Error())
		}

		usuario.Id = id
		usuario.Nome = nome
		usuario.Email = email
	}

	return usuario
}

// AtualizarUsuario faz o UPDATE no banco de dados
func AtualizarUsuario(id int, nome, email string) error {
	db := config.ConectarBD()
	defer db.Close()

	// Prepara o UPDATE
	atualizaDados, err := db.Prepare("UPDATE usuarios SET nome=$1, email=$2 WHERE id=$3")
	if err != nil {
		return err
	}
	defer atualizaDados.Close()

	// Executa passando os novos valores e o ID
	_, err = atualizaDados.Exec(nome, email, id)
	return err
}

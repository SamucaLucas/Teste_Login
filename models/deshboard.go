package models

import "app_login/config"

// BuscarTodosUsuarios retorna uma lista com todos os usuários do banco
func BuscarTodosUsuarios() []Usuario {
	db := config.ConectarBD()
	defer db.Close()

	// Seleciona os dados, ordenando pelo ID
	resultado, err := db.Query("SELECT id, nome, email FROM usuarios ORDER BY id ASC")
	if err != nil {
		panic(err.Error()) // Em produção, é melhor fazer um log e retornar o erro
	}
	defer resultado.Close()

	// Cria um slice (uma lista) vazia do tipo Usuario
	usuario := Usuario{}
	usuarios := []Usuario{}

	// Percorre cada linha que voltou do banco de dados
	for resultado.Next() {
		var id int
		var nome, email string

		// Escaneia os valores da linha atual para as variáveis
		err = resultado.Scan(&id, &nome, &email)
		if err != nil {
			panic(err.Error())
		}

		// Preenche a struct e adiciona na nossa lista
		usuario.Id = id
		usuario.Nome = nome
		usuario.Email = email

		usuarios = append(usuarios, usuario)
	}

	return usuarios
}

// DeletarUsuario remove um registro do banco de dados pelo ID
func DeletarUsuario(id string) {
	db := config.ConectarBD()
	defer db.Close()

	// Prepara a query de DELETE
	deletarRegistro, err := db.Prepare("DELETE FROM usuarios WHERE id=$1")
	if err != nil {
		panic(err.Error())
	}
	defer deletarRegistro.Close()

	// Executa a exclusão passando o ID
	_, err = deletarRegistro.Exec(id)
	if err != nil {
		panic(err.Error())
	}
}

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


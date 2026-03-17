package controllers

import (
	"app_login/models"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// Compila os templates HTML
var tmpl = template.Must(template.ParseGlob("views/*.html"))

var store = sessions.NewCookieStore([]byte("minha-chave-super-secreta-12345"))

func TelaLogin(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func TelaCadastro(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "cadastro.html", nil)
}

func ProcessarCadastro(w http.ResponseWriter, r *http.Request) {
	// Verifica se o método de envio é POST
	if r.Method == "POST" {
		// Pega os valores preenchidos no formulário (names dos inputs)
		nome := r.FormValue("nome")
		email := r.FormValue("email")
		senha := r.FormValue("senha")

		// Criptografa a senha gerando um Hash
		hashSenha, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Erro ao gerar hash da senha:", err)
			http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
			return
		}

		// Manda os dados para o Model salvar no PostgreSQL (convertendo o hash de byte para string)
		err = models.CriarUsuario(nome, email, string(hashSenha))
		if err != nil {
			http.Error(w, "Erro ao cadastrar. Verifique se o e-mail já está em uso.", http.StatusInternalServerError)
			return
		}

		// Se deu tudo certo, redireciona o usuário de volta para a tela de login
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func ProcessarLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		senha := r.FormValue("senha")

		usuario, err := models.ValidarLogin(email, senha)
		if err != nil {
			// A MUDANÇA ESTÁ AQUI 👇
			// Em vez de redirecionar (http.Redirect), renderizamos o template enviando a mensagem
			tmpl.ExecuteTemplate(w, "login.html", "E-mail ou senha incorretos!")
			return
		}

		// Recupera ou cria a sessão
		session, _ := store.Get(r, "sessao-app-vidro")

		// Define que o utilizador está autenticado e guarda o seu ID
		session.Values["autenticado"] = true
		session.Values["usuario_id"] = usuario.Id
		session.Save(r, w)

		// Entra no sistema
		http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
	}
}

// Logout destrói a sessão e sai do sistema
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessao-app-vidro")
	session.Values["autenticado"] = false
	session.Options.MaxAge = -1 // Apaga o cookie
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// --- ROTAS PROTEGIDAS ---

// Dashboard (agora com proteção)
func Dashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessao-app-vidro")

	// Verifica se a variável "autenticado" existe e é verdadeira
	if auth, ok := session.Values["autenticado"].(bool); !ok || !auth {
		// Se não estiver logado, expulsa para o login
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	todosOsUsuarios := models.BuscarTodosUsuarios()
	tmpl.ExecuteTemplate(w, "dashboard.html", todosOsUsuarios)
}

// Editar (agora com proteção)
func Editar(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessao-app-vidro")
	if auth, ok := session.Values["autenticado"].(bool); !ok || !auth {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	idDoUsuario := r.URL.Query().Get("id")
	usuario := models.BuscarUsuario(idDoUsuario)
	tmpl.ExecuteTemplate(w, "editar.html", usuario)
}

// Atualizar (agora com proteção)
func Atualizar(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessao-app-vidro")
	if auth, ok := session.Values["autenticado"].(bool); !ok || !auth {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	// Atualizar recebe o formulário de edição e salva no banco
	if r.Method == "POST" {
		id := r.FormValue("id")
		nome := r.FormValue("nome")
		email := r.FormValue("email")

		// Converte o ID de string para inteiro
		idConvertidoParaInt, err := strconv.Atoi(id)
		if err != nil {
			log.Println("Erro na conversão do ID para int:", err)
			http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
			return
		}

		// Manda atualizar no banco
		err = models.AtualizarUsuario(idConvertidoParaInt, nome, email)
		if err != nil {
			log.Println("Erro ao atualizar:", err)
		}
	}

	// Volta para o dashboard
	http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
}

// Deletar remove o usuário e volta para o dashboard
func Deletar(w http.ResponseWriter, r *http.Request) {
	// Proteção: verifica se está logado
	session, _ := store.Get(r, "sessao-app-vidro")
	if auth, ok := session.Values["autenticado"].(bool); !ok || !auth {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	// Pega o ID da URL
	idDoUsuario := r.URL.Query().Get("id")

	// Chama o model para deletar no banco
	models.DeletarUsuario(idDoUsuario)

	// Redireciona de volta para a lista atualizada
	http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
}

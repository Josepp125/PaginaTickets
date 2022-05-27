package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const urlMysql = "root:123654@tcp(localhost:3306)/base"

//Estructuras
func Saludar(nombre string) string {
	return "Hola " + nombre + " desde la función"
}

type Usuario struct {
	Id            int
	Mail          string
	Nombre        string
	NombreUsuario string
	Edad          int
	Contrasenia   string
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", urlMysql)
	db.Exec("create table if not exists usuario(id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,mail VARCHAR(50),nombre VARCHAR(30) NOT NULL,usuario VARCHAR(30) NOT NULL,edad INT NOT NULL,password VARCHAR(64) NOT NULL)")

	if err != nil {
		panic(err.Error())
	}
	log.Println("Base de datos conectada")
	return db
}

func Insert(rw http.ResponseWriter, r *http.Request) {
	usuario := Usuario{}
	db := dbConn()
	if r.Method == "POST" {
		mail := r.FormValue("mail")
		nombre := r.FormValue("nombre")
		usuario := r.FormValue("usuario")
		edad := r.FormValue("edad")
		password := r.FormValue("password")
		insForm, err := db.Prepare("INSERT INTO usuario(mail,nombre,usuario,edad,password) VALUES(?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(mail, nombre, usuario, edad, password)
		log.Println("Nuevo registro:" + mail + ", " + nombre + ", " + usuario + ", " + edad + ", " + password)
	}
	defer db.Close()
	renderTemplate(rw, "login.html", usuario)
}

// var templates = template.Must(template.New("T").ParseGlob("templates/*.html"))
var templates = template.Must(template.New("T").ParseGlob("templates/**/*.html"))

//función para renderizar los templates desde cada Handler
/**si se renderiza un archivo inexistente da "La conexión ha sido reiniciada"
* por lo que hay que manejar el error con http.Error() **/
var errorTemplate = template.Must(template.ParseFiles("templates/error/error.html"))

//HandlerError
func manejaError(rw http.ResponseWriter, status int) {
	rw.WriteHeader(status) //incluye el StatusError en el mensaje de eerror
	errorTemplate.Execute(rw, nil)
}

func renderTemplate(rw http.ResponseWriter, archivo string, data interface{}) {
	err := templates.ExecuteTemplate(rw, archivo, data)
	if err != nil {
		//http.Error(rw, "No es posible retornar template", http.StatusInternalServerError)
		manejaError(rw, http.StatusInternalServerError)
	}
}

//Handler
func Index(rw http.ResponseWriter, r *http.Request) {
	usuario := Usuario{}
	//renderTemplate(rw, "inde.html", usuario) //produce el error
	renderTemplate(rw, "index.html", usuario)
}

func Acercade(rw http.ResponseWriter, r *http.Request) {
	usuario := Usuario{}
	renderTemplate(rw, "acercade.html", usuario)
}

func Registrarse(rw http.ResponseWriter, r *http.Request) {
	usuario := Usuario{}
	renderTemplate(rw, "registro.html", usuario)
}
func Iniciarsesion(rw http.ResponseWriter, r *http.Request) {
	usuario := Usuario{}
	renderTemplate(rw, "login.html", usuario)
}
func Preguntas(rw http.ResponseWriter, r *http.Request) {
	usuario := Usuario{}
	renderTemplate(rw, "preguntas.html", usuario)
}
func Validar(rw http.ResponseWriter, r *http.Request) {
	nombre := r.URL.Query().Get("login")
	password := r.URL.Query().Get("pass")
	usuario := Usuario{}
	if nombre == "jose" && password == "123456" {
		renderTemplate(rw, "perfil.html", usuario)
	} else {
		renderTemplate(rw, "error1.html", usuario)
	}
}

func main() {

	//Conexion base de datos
	dbConn()
	//Archivos estáticos
	archEstaticos := http.FileServer(http.Dir("estaticos"))

	//Mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", Index)
	mux.HandleFunc("/acercade", Acercade)
	mux.HandleFunc("/preguntas", Preguntas)
	mux.HandleFunc("/registrarse", Registrarse)
	mux.HandleFunc("/insert", Insert)
	mux.HandleFunc("/login", Iniciarsesion)
	mux.HandleFunc("/validar", Validar)

	//Mux de archivos estáticos
	mux.Handle("/estaticos/", http.StripPrefix("/estaticos/", archEstaticos))

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	fmt.Println("Servidor corriendo en http://localhost:8080/")
	log.Fatal(server.ListenAndServe())

}

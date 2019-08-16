package main

import (
	//"fmt"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CezarGarrido/cuco_robots/api/driver"
	appHandler "github.com/CezarGarrido/cuco_robots/api/handler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "C102030g"
	dbname   = "bot_uems"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)

	fmt.Println("Local pg info:", psqlInfo)

	url, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatalln("$DATABASE_URL is required")
	}
	fmt.Println("Heroku pg info:", url)

	connection, err := driver.ConnectSQL(url)
	if err != nil {
		log.Panic(err)
	}
	alunoHandler := appHandler.NewAluno(connection)
	disciplinaHandler := appHandler.NewAlunoDisciplina(connection)
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/login", alunoHandler.Login).Methods("POST")
	r.HandleFunc("/api/v1/disciplinas", disciplinaHandler.Fetch).Methods("GET")

	headersOk := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})

	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE"})
	log.Println("Servidor startado na porta ", port)

	recoveryH := handlers.RecoveryHandler()(r)
	err = http.ListenAndServe(":"+port, handlers.CompressHandler(handlers.CORS(headersOk, methodsOk, originsOk)(recoveryH)))
	if err != nil {
		log.Panic(err)
	}
}

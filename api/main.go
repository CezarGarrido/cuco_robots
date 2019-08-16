package main

import (
	//"fmt"
	"log"
	"net/http"

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
	connection, err := driver.ConnectSQL(host, port, user, password, dbname)
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
	log.Println("Servidor startado na porta :8091")

	recoveryH := handlers.RecoveryHandler()(r)
	err = http.ListenAndServe(":8091", handlers.CompressHandler(handlers.CORS(headersOk, methodsOk, originsOk)(recoveryH)))
	if err != nil {
		log.Panic(err)
	}
}

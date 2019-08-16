package main

import (
	//"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CezarGarrido/cuco_robots/api/driver"
	appHandler "github.com/CezarGarrido/cuco_robots/api/handler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// const (
// 	host     = "localhost"
// 	port     = "5432"
// 	user     = "postgres"
// 	password = "C102030g"
// 	dbname   = "bot_uems"
// )
const (
	host     = "ec2-54-225-242-183.compute-1.amazonaws.com"
	port     = "5432"
	user     = "aimzpnysofwypw"
	password = "de56c756197c4d8f41745acf76ff3df6c3cc39852c7eb5572d173778d7ba28de"
	dbname   = "dbif64ksnitjje"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
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
	log.Println("Servidor startado na porta ", port)

	recoveryH := handlers.RecoveryHandler()(r)
	err = http.ListenAndServe(":"+port, handlers.CompressHandler(handlers.CORS(headersOk, methodsOk, originsOk)(recoveryH)))
	if err != nil {
		log.Panic(err)
	}
}

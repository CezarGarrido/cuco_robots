package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CezarGarrido/cuco_robots/api/driver"
	entities "github.com/CezarGarrido/cuco_robots/api/entities"
	repo "github.com/CezarGarrido/cuco_robots/api/repository"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("C_hello_uems_G")

func NewAluno(db *driver.DB) *Aluno {
	return &Aluno{
		repo: repo.NewSQLAlunoRepo(db.SQL),
	}
}

type Aluno struct {
	repo repo.AlunoRepo
}

func (p *Aluno) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Fprint(os.Stderr, "request cancelled\n")
	case <-time.After(2 * time.Second):
		creds := &entities.Aluno{}
		err := json.NewDecoder(r.Body).Decode(creds)
		if err != nil {
			log.Println(err.Error())
			respondWithError(w, 500, "Login ou senha inválidos")
			return
		}
		payload, err := p.repo.GetByLogin(ctx, creds.Rgm)
		if err != nil {
			log.Println(err.Error())
			respondWithError(w, 500, "Rgm inválido")
			return
		}
		if payload.Senha != creds.Senha {
			log.Println("Login ou senha inválidos", payload.Senha, creds.Senha)
			respondWithError(w, 500, "Senha inválida")
			return
		}
		//expirationTime := time.Now().Add(20 * time.Minute)
		claims := &entities.Claims{
			Aluno: *payload,
			/*StandardClaims: jwt.StandardClaims{
				//ExpiresAt: expirationTime.Unix(),
			},*/
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			log.Println(err.Error())
			respondWithError(w, 500, "Server Error")
			return
		}
		respondwithJSON(w, 200, tokenString)
	}
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}

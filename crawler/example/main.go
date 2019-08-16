package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/CezarGarrido/cuco_robots/crawler"
)

func main() {
	fmt.Println("# Fazendo login")
	client, err := crawler.NewClient("40089", "C102030g")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println("# Login Efetuado")
	fmt.Println("# Buscando dados do aluno")
	aluno, err := client.FindAluno()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	b, err := json.Marshal(aluno)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println(string(b))

	disciplinas, err := client.FindDisciplinas()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	bdisciplinas, err := json.Marshal(disciplinas)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println(string(bdisciplinas))
	for _, disciplina := range disciplinas {

		detalhe, _ := client.FindNotasByDisciplina(strconv.FormatInt(disciplina.UemsID, 10))
		bdetalhe, err := json.Marshal(detalhe)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		fmt.Println(string(bdetalhe))
	}
	_,_ = client.Logout()
}

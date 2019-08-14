package main

import (
	"encoding/json"
	"fmt"

	"github.com/CezarGarrido/cuco_robots/crawler"
)

func main() {
	client, err := crawler.NewClient("40089", "C102030g")
	if err != nil {
		panic(err)
	}
	aluno, err := client.FindAluno()
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal(aluno)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println(string(b))

}

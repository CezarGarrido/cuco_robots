package main

import (
	"fmt"

	"github.com/CezarGarrido/cuco_robots/crawler"
)

func main() {
	client, err := crawler.NewClient("rgm", "senha")
	if err != nil {
		panic(err)
	}
	aluno, err := client.FindAluno()
	if err != nil {
		panic(err)
	}
	fmt.Println(*aluno)

}

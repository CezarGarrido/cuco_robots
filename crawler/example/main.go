package main

import (
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
	fmt.Println(*aluno)

}

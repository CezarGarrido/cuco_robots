package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/CezarGarrido/cuco_robots/crawler"
)

func main() {
	start := time.Now()
	fmt.Println("# Fazendo login")
	client, err := crawler.NewClient("40089", "C102030g")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println("# Login Efetuado")
	fmt.Println("# Buscando dados do aluno")
	falta, err := client.FindFaltas("1391920")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	b, err := json.Marshal(falta)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println(string(b))
	/*aluno, err := client.FindAluno()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	b, err := json.Marshal(aluno)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println(string(b))*/
	//for _, c := range client.GetCookies() {
	//	fmt.Println(c.Name, c.Value, c.Path, c.Domain)
	//	}
	//reqCookie(client.GetCookies())

	/*disciplinas, err := client.FindDisciplinas()
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
	}*/
	//_, _ = client.Logout()
	fmt.Println(time.Since(start))
}

func reqCookie(Cookies []*http.Cookie) {
	req, err := http.NewRequest("GET", "https://sistemas.uems.br/academico/dcu005.php", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	HtppClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	for _, cookie := range Cookies {
		req.AddCookie(cookie)
	}
	resposta, err := HtppClient.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	defer resposta.Body.Close()
	body, err := ioutil.ReadAll(resposta.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println(string(body))

}

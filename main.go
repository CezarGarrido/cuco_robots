package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"./internal"
	"./models"
)
func main() {
	db, err := internal.ConexaoPostgres()
	if err != nil {
		panic(err)
	}
	aluno := models.Aluno{}
	alunos, err := aluno.GetAll(db)
	if err != nil {
		panic(err)
	}
	for _, auxAluno := range alunos {
		consultarNotas(auxAluno)
	}
}
func consultarNotas(aluno *models.Aluno) {
	parm := url.Values{}
	parm.Add("login", "")
	parm.Add("rgm", aluno.Rgm)
	parm.Add("senha", aluno.Senha)
	cookieJar, _ := cookiejar.New(nil)
	req, err := http.NewRequest("POST", "https://sistemas.uems.br/academico/index.php", strings.NewReader(parm.Encode()))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{
		Jar: cookieJar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer resp.Body.Close()
	parm2 := url.Values{}
	parm2.Add("event", "notas")
	parm2.Add("list[matricula_aluno_turma.codigo]", "1391920")
	req, err = http.NewRequest("POST", "https://sistemas.uems.br/academico/dcu003.php", strings.NewReader(parm2.Encode()))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp1, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer resp1.Body.Close()
	body1, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}
	db, err := internal.ConexaoPostgres()
	if err != nil {
		panic(err)
	}
	nota := models.Nota{}
	nota.IdAluno = aluno.Id
	nota.IdDisciplina = 1334
	nota.Documento = string(body1)
	err = nota.Create(db)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}




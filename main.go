package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/CezarGarrido/cuco_robots/driver"
	"github.com/CezarGarrido/cuco_robots/models"
	"github.com/PuerkitoBio/goquery"
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
		client, err := newClient(auxAluno)
		if err != nil {
			fmt.Printf("%+v\n", err.Error())
			return
		}
		disciplinas, err := consultarDisciplinas(auxAluno, client)
		if err != nil {
			fmt.Printf("%+v\n", err.Error())
			return
		}
		for _, disciplina := range disciplinas {
			fmt.Println(disciplina.ID)
			uemsDisciplina := models.Disciplina{}
			uemsDisciplina.IDUEMS = disciplina.ID
			uemsDisciplina.Descricao = disciplina.Descricao
			uemsDisciplina.CreatedAt = time.Now()
			if !uemsDisciplina.IsExist(db) {
				idDisc, err := uemsDisciplina.Create(db)
				if err != nil {
					return
				}
				uemsDisciplina.ID = idDisc
			} else {
				err = uemsDisciplina.GetByIDUEMS(uemsDisciplina.IDUEMS, db)
				if err != nil {
					return
				}
				fmt.Println("> Disciplina já cadastrada no sistema")
			}
			alunoDisciplina := models.AlunoDisciplina{}
			alunoDisciplina.IDAluno = auxAluno.ID
			alunoDisciplina.IDDisciplina = uemsDisciplina.ID
			alunoDisciplina.IDUEMS = uemsDisciplina.IDUEMS
			alunoDisciplina.CreatedAt = time.Now()
			if !alunoDisciplina.IsExist(db) {
				fmt.Println("> Cadastrando disciplina para o aluno")
				_, err = alunoDisciplina.Create(db)
				if err != nil {
					return
				}

			} else {
				fmt.Println("> Disciplina já cadastrada para o aluno")
			}
			nota := models.Nota{}

			nota.IDAluno = auxAluno.ID
			nota.IDDisciplina = uemsDisciplina.ID
			strID := strconv.FormatInt(alunoDisciplina.IDUEMS, 10)
			doc, err := consultarNotas(strID, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			nota.Documento = *doc
			if nota.IsExist(db) {
				err = nota.GetByAluno(nota.IDAluno, db)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				err := nota.Update(db)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			} else {
				_, err = nota.Create(db)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		}
	}
}

func newClient(aluno *models.Aluno) (*http.Client, error) {
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
		return nil, err
	}
	defer resp.Body.Close()

	return client, nil
}

func consultarNotas(idDisciplina string, client *http.Client) (*string, error) {

	param := url.Values{}
	param.Add("event", "notas")
	param.Add("list[matricula_aluno_turma.codigo]", idDisciplina)
	req, err := http.NewRequest("POST", "https://sistemas.uems.br/academico/dcu003.php", strings.NewReader(param.Encode()))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc := string(body)
	return &doc, nil
}

func consultarDisciplinas(aluno *models.Aluno, client *http.Client) ([]Disciplina, error) {
	req, err := http.NewRequest("POST", "https://sistemas.uems.br/academico/dcu003.php", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp1, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()
	body1, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return nil, err
	}
	disciplinas, err := parserDisciplinas(string(body1))
	if err != nil {
		return nil, err
	}
	return disciplinas, nil
}

type Disciplina struct {
	ID        int64
	Descricao string
	Oferta    string
}

func parserDisciplinas(html string) ([]Disciplina, error) {

	disciplinas := []Disciplina{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return disciplinas, err
	}
	//doc.Find("script").Remove()
	//doc.Find("style").Remove()
	//doc.Find("img").Remove()
	var disciplina Disciplina
	isError := false
	var erros []string
	doc.Find("table.event_list").Each(func(index int, tablehtml *goquery.Selection) {
		//band, ok := tablehtml.Attr("class")
		//if ok {
		//if band == "event_list" {

		tablehtml.Find("tr#link").Each(func(indextr int, rowhtml *goquery.Selection) {
			band, ok := rowhtml.Attr("onclick")
			if ok {
				re := regexp.MustCompile("[0-9]+")
				id := strings.Join(re.FindAllString(band, -1), "")
				n, err := strconv.ParseInt(id, 10, 64)
				if err != nil {
					erros = append(erros, err.Error())
					isError = true
				}
				disciplina.ID = n
				rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
					if indexth == 0 {
						disciplina.Descricao = strings.Join(strings.Fields(tablecell.Text()), " ")
					} else if indexth == 1 {
						disciplina.Oferta = strings.Join(strings.Fields(tablecell.Text()), " ")
					}
				})
			}
			disciplinas = append(disciplinas, disciplina)
		})
		//}
		//}
	})
	if isError {
		return disciplinas, errors.New(strings.Join(erros, " "))
	}
	return disciplinas, nil
}

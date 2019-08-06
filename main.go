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

// Disciplina formatos usando os formatos predefinidos para seus operandos e escreve para w.
// Espaços são adicionados entre os operandos quando nem é uma seqüência de caracteres.
// Retorna o número de bytes gravados e qualquer erro de escrita encontrados.
type Disciplina struct {
	ID        int64
	Descricao string
	Oferta    string
}

// Detalhes formatos usando os formatos predefinidos para seus operandos e escreve para w.
// Espaços são adicionados entre os operandos quando nem é uma seqüência de caracteres.
// Retorna o número de bytes gravados e qualquer erro de escrita encontrados.
type Detalhes struct {
	Unidade                string
	Curso                  string
	Disciplina             string
	Turma                  string
	SerieDisciplina        string
	CargaHorariaPresencial string
	MaximoFaltas           string
	PeriodoLetivo          string
	Professor              string
	MediaAvaliacoes        string
	Optativa               string
	Exame                  string
	MediaFinal             string
	Faltas                 string
	Situacao               string
	Notas                  []Nota
}

// Nota formatos usando os formatos predefinidos para seus operandos e escreve para w.
// Espaços são adicionados entre os operandos quando nem é uma seqüência de caracteres.
// Retorna o número de bytes gravados e qualquer erro de escrita encontrados.
type Nota struct {
	Descricao string
	Valor     string
}

// Start formatos usando os formatos predefinidos para seus operandos e escreve para w.
// Espaços são adicionados entre os operandos quando nem é uma seqüência de caracteres.
// Retorna o número de bytes gravados e qualquer erro de escrita encontrados.
func Start() {
	fmt.Println("> 🔥 Starting")

	db, err := internal.ConexaoPostgres()
	if err != nil {
		panic(err)
	}
	aluno := models.Aluno{}
	fmt.Println("> buscando alunos")
	alunos, err := aluno.GetAll(db)
	if err != nil {
		panic(err)
	}
	for _, auxAluno := range alunos {
		fmt.Println("> iniciando sessão")
		client, err := newClient(auxAluno)
		if err != nil {
			fmt.Printf("%+v\n", err.Error())
			return
		}
		fmt.Println("> consultando disciplinas")
		disciplinas, err := consultarDisciplinas(auxAluno, client)
		if err != nil {
			fmt.Printf("%+v\n", err.Error())
			return
		}
		for _, disciplina := range disciplinas {
			uemsDisciplina := models.Disciplina{}
			uemsDisciplina.IDUEMS = disciplina.ID
			uemsDisciplina.Descricao = disciplina.Descricao
			uemsDisciplina.CreatedAt = time.Now()
			fmt.Println("> Disciplina", disciplina.Descricao)
			if !uemsDisciplina.IsExist(db) {
				idDisc, err := uemsDisciplina.Create(db)
				if err != nil {
					fmt.Println("Erro ao criar disciplina", err.Error())
					return
				}
				uemsDisciplina.ID = idDisc
			} else {
				err = uemsDisciplina.GetByIDUEMS(uemsDisciplina.IDUEMS, db)
				if err != nil {
					fmt.Println("Erro ao buscar disciplina", err.Error())
					return
				}
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
					fmt.Println("Erro ao cadastrar disciplina para o aluno", err.Error())
					return
				}

			}
			strID := strconv.FormatInt(alunoDisciplina.IDUEMS, 10)
			fmt.Println("> consultando notas")
			doc, err := consultarNotas(strID, client)
			if err != nil {
				fmt.Println("Erro ao consultar notas", err.Error())
				return
			}
			detalhe, err := parserNotas(*doc)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			alunoNota := models.AlunoNota{}
			alunoNota.IDAluno = auxAluno.ID
			alunoNota.IDDisciplina = uemsDisciplina.ID
			alunoNota.Unidade = detalhe.Unidade
			alunoNota.Curso = detalhe.Curso
			alunoNota.Disciplina = detalhe.Disciplina
			alunoNota.Turma = detalhe.Turma
			alunoNota.SerieDisciplina = detalhe.SerieDisciplina
			alunoNota.CargaHorariaPresencial = detalhe.CargaHorariaPresencial
			alunoNota.MaximoFaltas = detalhe.MaximoFaltas
			alunoNota.PeriodoLetivo = detalhe.PeriodoLetivo
			alunoNota.Professor = detalhe.Professor
			alunoNota.MediaAvaliacoes = detalhe.MediaAvaliacoes
			alunoNota.Optativa = detalhe.Optativa
			alunoNota.Exame = detalhe.Exame
			alunoNota.MediaFinal = detalhe.MediaFinal
			alunoNota.Faltas = detalhe.Faltas
			alunoNota.Situacao = detalhe.Situacao
			alunoNota.CreatedAt = time.Now()
			alunoNota.UpdatedAt = time.Now()
			if !alunoNota.IsExist(db) {
				notaID, err := alunoNota.Create(db)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				alunoNota.ID = notaID
			} else {
				err = alunoNota.GetByDisciplina(uemsDisciplina.ID, db)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				err = alunoNota.Update(db)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
			for _, n := range detalhe.Notas {
				fmt.Println(">", n.Descricao, " -->", n.Valor)
				valorNota := models.ValorNota{
					IDNota:    alunoNota.ID,
					Descricao: n.Descricao,
					Valor:     n.Valor,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				if !valorNota.IsExist(db) {
					_, err := valorNota.Create(db)
					if err != nil {
						fmt.Println("Erro ao criar valor nota", err.Error())
						return
					}
				} else {
					err = valorNota.GetByDescricao(n.Descricao, db)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					fmt.Println("> Atualizando nota")
					valorNota.UpdatedAt = time.Now()
					err = valorNota.Update(db)
					if err != nil {
						fmt.Println("erro ao atualizar valor da nota", err.Error())
						return
					}
				}
			}

		}
	}
}

// newClient formatos usando os formatos predefinidos para seus operandos e escreve para w.
// Espaços são adicionados entre os operandos quando nem é uma seqüência de caracteres.
// Retorna o número de bytes gravados e qualquer erro de escrita encontrados.
func newClient(aluno *models.Aluno) (*http.Client, error) {

	param := url.Values{}
	param.Add("login", "")
	param.Add("rgm", aluno.Rgm)
	param.Add("senha", aluno.Senha)

	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://sistemas.uems.br/academico/index.php", strings.NewReader(param.Encode()))
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

// consultarNotas formatos usando os formatos predefinidos para seus operandos e escreve para w.
// Espaços são adicionados entre os operandos quando nem é uma seqüência de caracteres.
// Retorna o número de bytes gravados e qualquer erro de escrita encontrados.
func consultarNotas(idDisciplina string, client *http.Client) (*string, error) {

	param := url.Values{}
	param.Add("event", "notas")
	param.Add("list[matricula_aluno_turma.codigo]", idDisciplina)
	req, err := http.NewRequest("POST", "https://sistemas.uems.br/academico/dcu003.php", strings.NewReader(param.Encode()))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
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

func parserDisciplinas(html string) ([]Disciplina, error) {

	disciplinas := []Disciplina{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return disciplinas, err
	}
	var disciplina Disciplina
	isError := false
	var erros []string
	doc.Find("table.event_list").Each(func(index int, tablehtml *goquery.Selection) {
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
	})
	if isError {
		return disciplinas, errors.New(strings.Join(erros, " "))
	}
	return disciplinas, nil
}

func parserNotas(html string) (Detalhes, error) {
	var detalhe Detalhes
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return detalhe, err
	}
	var nota Nota
	var notas []Nota
	doc.Find("table.event_form").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			band, ok := rowhtml.Attr("id")
			if ok {
				if band == "tr_unidade" {
					rowhtml.Find("td").Each(func(indexth int, tdtml *goquery.Selection) {
						detalhe.Unidade = strings.Join(strings.Fields(tdtml.Text()), " ")
					})
				}
				if band == "tr_curso" {
					rowhtml.Find("td").Each(func(indexth int, tdtml *goquery.Selection) {
						detalhe.Curso = strings.Join(strings.Fields(tdtml.Text()), " ")
					})
				}
				if band == "tr_disciplina" {
					rowhtml.Find("td").Each(func(indexth int, tdtml *goquery.Selection) {
						detalhe.Disciplina = strings.Join(strings.Fields(tdtml.Text()), " ")
					})
				}
				if band == "tr_turma" {
					rowhtml.Find("td").Each(func(indexth int, tdtml *goquery.Selection) {
						detalhe.Turma = strings.Join(strings.Fields(tdtml.Text()), " ")
					})
				}
				if band == "tr_serie" {
					rowhtml.Find("td").Each(func(indexth int, tdtml *goquery.Selection) {
						detalhe.SerieDisciplina = strings.Join(strings.Fields(tdtml.Text()), " ")
					})
				}
				if band == "tr_carga_horaria" {
					rowhtml.Find("td").Each(func(indexth int, tdtml *goquery.Selection) {
						detalhe.CargaHorariaPresencial = strings.Join(strings.Fields(tdtml.Text()), " ")
					})
				}
				//tr_curso
				//tr_disciplina
				//tr_turma
				//tr_serie
				//tr_carga_horaria
			} else {
				var texto string
				rowhtml.Find("th").Each(func(indexth int, tablecell *goquery.Selection) {
					texto = strings.Join(strings.Fields(tablecell.Text()), " ")
				})
				if texto == "Máximo de Faltas" {
					rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
						detalhe.MaximoFaltas = strings.Join(strings.Fields(tablecell.Text()), " ")
					})
				}
				if texto == "Período Letivo" {
					rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
						detalhe.PeriodoLetivo = strings.Join(strings.Fields(tablecell.Text()), " ")
					})
				}
				if texto == "Professor(a)" {
					rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
						detalhe.Professor = strings.Join(strings.Fields(tablecell.Text()), " ")
					})
				}
			}
			rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {

				band, ok := tablecell.Attr("colspan")
				if ok {
					if band == "2" {
						tablecell.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
							if indextr == 1 {
								rowhtml.Find("th").Each(func(indexth int, thtml *goquery.Selection) {
									nota.Descricao = strings.Join(strings.Fields(thtml.Text()), " ")
									notas = append(notas, nota)
								})
							}
							if indextr == 2 {
								rowhtml.Find("td").Each(func(indexth int, thtml *goquery.Selection) {
									if indexth < len(notas) {
										notas[indexth].Valor = strings.Join(strings.Fields(thtml.Text()), " ")
									} else {
										if indexth == len(notas) {
											fmt.Println(strings.Join(strings.Fields(thtml.Text()), " "))
											detalhe.MediaAvaliacoes = strings.Join(strings.Fields(thtml.Text()), " ")
										}
										if indexth == len(notas)+1 {
											detalhe.Optativa = strings.Join(strings.Fields(thtml.Text()), " ")
										}
										if indexth == len(notas)+2 {
											detalhe.Exame = strings.Join(strings.Fields(thtml.Text()), " ")
										}
										if indexth == len(notas)+3 {
											detalhe.MediaFinal = strings.Join(strings.Fields(thtml.Text()), " ")
										}
										if indexth == len(notas)+4 {
											detalhe.Faltas = strings.Join(strings.Fields(thtml.Text()), " ")
										}
										if indexth == len(notas)+5 {
											detalhe.Situacao = strings.Join(strings.Fields(thtml.Text()), " ")
										}
									}
								})
							}
						})
					}
				}
			})
		})
	})
	detalhe.Notas = notas
	return detalhe, nil
}

func main() {
	Start()
}

type Client struct {
	Conn *http.Client
}

func NewClient(rgm, senha string) (*Client, error) {
	
	client := new(Client)

	param := url.Values{}

	param.Add("login", "")
	param.Add("rgm", rgm)
	param.Add("senha", senha)

	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://sistemas.uems.br/academico/index.php", strings.NewReader(param.Encode()))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	HtppClient := &http.Client{
		Jar: cookieJar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := HtppClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	client.Conn = HtppClient
	return client, nil
}

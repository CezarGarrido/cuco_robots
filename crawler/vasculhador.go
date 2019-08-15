package crawler

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	BaseURL string
	Conn    *http.Client
}
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
type Aluno struct {
	Nome            string
	DataNascimento  time.Time
	Sexo            string
	NomePai         string
	NomeMae         string
	EstadoCivil     string
	Nacionalidade   string
	Naturalidade    string
	Fenotipo        string
	CPF             string
	RG              string
	RGOrgaoEmissor  string
	RGEstadoEmissor string
	RGDataEmissao   time.Time
	Contatos        []*Contato
	Enderecos       []*Endereco
}
type Contato struct {
	Tipo  string
	Valor string
}
type Endereco struct {
	Logradouro  string
	Numero      int
	Complemento string
	Bairro      string
	CEP         string
	Cidade      string
}
type Disciplina struct {
	ID        int64
	Descricao string
	Oferta    string
}
type Nota struct {
	Descricao string
	Valor     string
}

func NewClient(rgm, senha string) (Client, error) {
	client := Client{}
	client.BaseURL = "https://sistemas.uems.br/academico/index.php"
	param := url.Values{}
	param.Add("login", "")
	param.Add("rgm", rgm)
	param.Add("senha", senha)
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return client, err
	}
	req, err := http.NewRequest("POST", client.BaseURL, strings.NewReader(param.Encode()))
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
		return client, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return client, err
	}
	msg, ok := checkLoginError(string(body))
	if !ok {
		return client, errors.New(msg)
	}
	client.Conn = HtppClient
	return client, nil
}
func (c Client) GetCookies() []*http.Cookie {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	return c.Conn.Jar.Cookies(u)
}
func (c Client) Logout() (Client, error) {
	param := url.Values{}
	param.Add("acao", "fechar")
	req, err := http.NewRequest("POST", c.BaseURL, strings.NewReader(param.Encode()))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Conn.Do(req)
	if err != nil {
		return c, err
	}
	defer resp.Body.Close()
	return c, nil
}

func (c Client) FindAluno() (*Aluno, error) {
	req, err := http.NewRequest("GET", "https://sistemas.uems.br/academico/dcu005.php", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resposta, err := c.Conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resposta.Body.Close()
	body, err := ioutil.ReadAll(resposta.Body)
	if err != nil {
		return nil, err
	}
	return parserAluno(string(body))
}

//parserAluno: Função que pega os dados do aluno
func parserAluno(html string) (*Aluno, error) {
	aluno := &Aluno{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return aluno, err
	}
	doc.Find("body.uc table.event_form tbody tr").Each(func(index int, tablehtml *goquery.Selection) {
		if index <= 26 {
			tablehtml.Find("th").Each(func(indexth int, thhtml *goquery.Selection) {
				var tipo string
				var valor string
				tipo = strings.Join(strings.Fields(thhtml.Text()), " ")
				tablehtml.Find("td").Each(func(indextd int, tdhtml *goquery.Selection) {
					valor = strings.Join(strings.Fields(tdhtml.Text()), " ")
				})
				//fmt.Println(tipo)
				switch tipo {
				case "Nome":
					aluno.Nome = valor
				case "Data de Nascimento":
					aluno.DataNascimento, _ = time.Parse("02/01/2006", valor)
				case "Sexo":
					aluno.Sexo = valor
				case "Nome do Pai":
					aluno.NomePai = valor
				case "Nome da Mãe":
					aluno.NomeMae = valor
				case "Estado Civil":
					aluno.EstadoCivil = valor
				case "Nacionalidade":
					aluno.Nacionalidade = valor
				case "Naturalidade":
					aluno.Naturalidade = valor
				case "CPF":
					aluno.CPF = valor
				case "RG":
					aluno.RG = valor
				case "Órgão Emissor":
					aluno.RGOrgaoEmissor = valor
				case "Estado":
					aluno.RGEstadoEmissor = valor
				case "Data de Emissão":
					aluno.RGDataEmissao, _ = time.Parse("02/01/2006", valor)
				}
			})

		}
	})

	contatos, _ := parserContatosAluno(html)
	aluno.Contatos = contatos
	enderecos, _ := parserEnderecosAluno(html)
	aluno.Enderecos = enderecos

	//fmt.Println(contatos)

	return aluno, nil
}

//parserContatosAluno: função que pega os dados de contato do aluno
func parserContatosAluno(html string) ([]*Contato, error) {
	contatos := make([]*Contato, 0)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return contatos, err
	}
	doc.Find("tbody#SubDatasetField1_tbody").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indexth int, trhtml *goquery.Selection) {
			trhtml.Find("td").Each(func(indexth int, tdhtml *goquery.Selection) {
				tdhtml.Find("input").Each(func(indexth int, inputhtml *goquery.Selection) {
					band, ok := inputhtml.Attr("value")
					if ok {
						if band != "f" && band != "t" {
							contato := &Contato{}
							if isPhoneNumber(strings.Join(strings.Fields(band), " ")) {
								contato.Tipo = "Telefone"
								contato.Valor = strings.Join(strings.Fields(band), " ")
								contatos = append(contatos, contato)
							} else {
								contato.Tipo = "Email"
								contato.Valor = strings.Join(strings.Fields(band), " ")
								contatos = append(contatos, contato)
							}
						}
					}
				})
			})
		})
	})
	return contatos, nil
}

//parserEnderecosAluno: função que busca os endereços dos alunos cadastrados
func parserEnderecosAluno(html string) ([]*Endereco, error) {
	enderecos := make([]*Endereco, 0)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return enderecos, err
	}
	doc.Find("table#SubDatasetField2").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tbody#SubDatasetField2_tbody").Each(func(indextbody int, trhtml *goquery.Selection) {
			trhtml.Find("tr").Each(func(indextr int, thhtml *goquery.Selection) {
				endereco := &Endereco{}
				thhtml.Find("td").Each(func(indextd int, tdhtml *goquery.Selection) {
					if indextd == 5 {
						endereco.Cidade = buscaCidade(indextd, tdhtml)
					}
					tdhtml.Find("input").Each(func(indexinput int, inputhtml *goquery.Selection) {
						band, ok := inputhtml.Attr("value")
						if ok {
							if band != "f" && band != "t" {
								if indextd == 0 {
									endereco.Logradouro = strings.Join(strings.Fields(band), " ")
								}
								if indextd == 1 {
									endereco.Numero, _ = strconv.Atoi(strings.Join(strings.Fields(band), " "))
								}
								if indextd == 2 {
									endereco.Complemento = strings.Join(strings.Fields(band), " ")
								}
								if indextd == 3 {
									endereco.Bairro = strings.Join(strings.Fields(band), " ")
								}
								if indextd == 4 {
									endereco.CEP = strings.Join(strings.Fields(band), " ")
								}
							}
						}

					})
				})
				enderecos = append(enderecos, endereco)
			})
		})
	})
	return enderecos, nil
}
func buscaCidade(index int, tdhtml *goquery.Selection) string {
	var cidade string
	tdhtml.Find("select").Each(func(indexselect int, selecthtml *goquery.Selection) {
		selecthtml.Find("option").Each(func(indexopt int, optionhtml *goquery.Selection) {
			_, ok := optionhtml.Attr("selected")
			if ok {
				cidade = strings.Join(strings.Fields(optionhtml.Text()), " ")
				return
			}
		})
	})
	return cidade
}
func (c Client) FindNotasByDisciplina(idDisciplina string) (Detalhes, error) {
	var detalhe Detalhes
	param := url.Values{}
	param.Add("event", "notas")
	param.Add("list[matricula_aluno_turma.codigo]", idDisciplina)
	req, err := http.NewRequest("POST", "https://sistemas.uems.br/academico/dcu003.php", strings.NewReader(param.Encode()))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Conn.Do(req)
	if err != nil {
		return detalhe, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return detalhe, err
	}
	defer resp.Body.Close()
	//doc := string(body)
	return parserNotas(string(body))
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

func (c Client) FindDisciplinas() ([]*Disciplina, error) {
	req, err := http.NewRequest("POST", "https://sistemas.uems.br/academico/dcu003.php", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parserDisciplinas(string(body))
}

func parserDisciplinas(html string) ([]*Disciplina, error) {
	disciplinas := make([]*Disciplina, 0)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return disciplinas, err
	}

	isError := false
	var erros []string
	doc.Find("table.event_list").Each(func(index int, tablehtml *goquery.Selection) {

		tablehtml.Find("tr#link").Each(func(indextr int, rowhtml *goquery.Selection) {
			var disciplina Disciplina
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
			disciplinas = append(disciplinas, &disciplina)
		})
	})
	if isError {
		return disciplinas, errors.New(strings.Join(erros, " "))
	}
	return disciplinas, nil
}

//checkLoginError: Função que recebe o html retornado na pagina de login
//e checa se existe a class error no html,
//se a class existir retorna o texto que esta na class e false,
//se não, retorna a mensagem de Bem-vindo e true
func checkLoginError(html string) (string, bool) {
	var isLoggged = true
	var msg = "Bem-vindo"
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "Erro ao parsear html de login", false
	}
	doc.Find(".error").Each(func(index int, errorhtml *goquery.Selection) {
		msg = strings.Join(strings.Fields(errorhtml.Text()), " ")
		isLoggged = false
	})
	return msg, isLoggged
}

//isPhoneNumber: verifica se uma string é um numero de telefone valido
func isPhoneNumber(number string) (ok bool) {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	submatch := re.FindStringSubmatch(number)
	ok = true
	if len(submatch) < 2 {
		ok = false
	}
	return
}

package crawler

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
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
	Conn *http.Client
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

//Logradouro	Nº	Complemento	Bairro	CEP	Cidade
func NewClient(rgm, senha string) (Client, error) {
	client := Client{}
	param := url.Values{}
	param.Add("login", "")
	param.Add("rgm", rgm)
	param.Add("senha", senha)
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return client, err
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
	//fmt.Println(contatos)
	enderecos, _ := parserEnderecosAluno(html)
	aluno.Enderecos = enderecos
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
					tdhtml.Find("select").Each(func(indexselect int, selecthtml *goquery.Selection) {
						selecthtml.Find("option").Each(func(indexopt int, optionhtml *goquery.Selection) {
							_, ok := optionhtml.Attr("selected")
							if ok {
								endereco.Cidade = strings.Join(strings.Fields(optionhtml.Text()), " ")
							}
						})
					})
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

func isPhoneNumber(number string) (ok bool) {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	submatch := re.FindStringSubmatch(number)
	ok = true
	if len(submatch) < 2 {
		ok = false
	}
	return
}

/*
# th: Nome
# td: CEZAR GARRIDO BRITEZ
# th: Data de Nascimento
# td: 28/12/1997
# th: Sexo
# td: Masculino
# th: Nome do Pai
# td: VITOR BRITEZ
# th: Nome da Mãe
# td: MARIANA GARRIDO
# th: Estado Civil
# td: Solteiro(a)
# th: Nacionalidade
# td: BRASILEIRO
# th: Naturalidade
# td: PARANHOS/MS
# th: Fenótipo *
# td: AmarelaBrancaIndígenaNão declaradoPardaPreta
# th: CPF
# td: 050.433.691-67
# th: DOCUMENTO DE IDENTIFICAÇÃO
# th: RG
# td: 2.225.228
# th: Órgão Emissor
# td: SEJUSP
# th: Estado
# td: MS
# th: Data de Emissão
# td:
# th: TÍTULO ELEITORAL
# th: Número
# td: 027682941937
# th: Zona
# td: 43
# th: Seção
# td: 242
# th: Cidade
# td: DOURADOS/MS
# th: Data de Expedição
# td: 08/03/2018
# th: ALISTAMENTO MILITAR
# th: Reservista
# td: 32000058497-3
# th: Categoria
# td:
# th: Órgão Emissor
# td: MD
# th: Série
# td:
# th: Data de Expedição
# td: 30/01/2018
*/

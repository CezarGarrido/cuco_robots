package crawler

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	Conn *http.Client
}

type Aluno struct {
	Nome string
}
type Contato struct {
	Tipo  string
	Valor string
}
type Endereco struct {

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
				//fmt.Println("# th:", strings.Join(strings.Fields(thhtml.Text()), " "))
			})
			tablehtml.Find("td").Each(func(indextd int, tdhtml *goquery.Selection) {
				//fmt.Println("# td:", strings.Join(strings.Fields(tdhtml.Text()), " "))
			})
		}
	})
	//parserContatosAluno(html)
	parserEnderecosAluno(html)
	return aluno, nil
}

//parserContatosAluno: função que pega os dados de contato do aluno
func parserContatosAluno(html string) (*Contato, error) {
	contato := &Contato{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return contato, err
	}
	doc.Find("tbody#SubDatasetField1_tbody").Each(func(index int, tablehtml *goquery.Selection) {
		fmt.Println("# SubDatasetField1_tbody:", strings.Join(strings.Fields(tablehtml.Text()), " "))
		tablehtml.Find("tr").Each(func(indexth int, trhtml *goquery.Selection) {
			trhtml.Find("td").Each(func(indexth int, tdhtml *goquery.Selection) {
				tdhtml.Find("input").Each(func(indexth int, inputhtml *goquery.Selection) {
					band, ok := inputhtml.Attr("value")
					if ok {
						if band != "f" && band != "t" {
							fmt.Println("# td input:", strings.Join(strings.Fields(band), " "))
							fmt.Println("is Phone:", isPhoneNumber(strings.Join(strings.Fields(band), " ")))
						}
					}
				})
			})
		})
	})
	return contato, nil
}

func parserEnderecosAluno(html string)(*Endereco, error){
	endereco:= &Endereco{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return endereco, err
	}
	doc.Find("table#SubDatasetField2").Each(func(index int, tablehtml *goquery.Selection) {
		//fmt.Println("# SubDatasetField2_tbody:", strings.Join(strings.Fields(tablehtml.Text()), " "))
		tablehtml.Find("tr").Each(func(indextr int, trhtml *goquery.Selection) {
			trhtml.Find("th").Each(func(indexth int, thhtml *goquery.Selection) {
				fmt.Println("# th:", strings.Join(strings.Fields(thhtml.Text()), " "),"indice tr:", indextr, "indice th:", indexth)
			})
		})
		tablehtml.Find("tbody#SubDatasetField2_tbody").Each(func(indextbody int, trhtml *goquery.Selection) {
			trhtml.Find("tr").Each(func(indextr int, thhtml *goquery.Selection) {
				thhtml.Find("td").Each(func(indextd int, tdhtml *goquery.Selection) {
					tdhtml.Find("input").Each(func(indexinput int, inputhtml *goquery.Selection) {
						band, ok := inputhtml.Attr("value")
						if ok {
							if band != "f" && band != "t" {
								fmt.Println("# td input:", strings.Join(strings.Fields(band), " "), "indice:", indextd)
								fmt.Println("is Phone:", isPhoneNumber(strings.Join(strings.Fields(band), " ")))
							}
						}
					})
					
				})
			})
		})
	})
	return endereco, nil
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

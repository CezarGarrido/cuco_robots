package crawler

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	Conn *http.Client
}

type Aluno struct {
	Nome string
}

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

func parserAluno(html string) (*Aluno, error) {
	aluno := &Aluno{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return aluno, err
	}
	doc.Find("table#table_event_form").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			band, ok := rowhtml.Attr("id")
			fmt.Println("# ok:", ok, "band:", band)
			if !ok {
				rowhtml.Find("th").Each(func(indexth int, thhtml *goquery.Selection) {
					fmt.Println(strings.Join(strings.Fields(thhtml.Text()), " "))
				})
				rowhtml.Find("td").Each(func(indextd int, tdhtml *goquery.Selection) {
					fmt.Println(strings.Join(strings.Fields(tdhtml.Text()), " "))
				})
			}

		})
	})
	return aluno, nil
}

func checkLoginError(html string) (string, bool) {
	var isLoggged = true
	var msg = "Bem-vindo"
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", false
	}
	doc.Find(".error").Each(func(index int, errorhtml *goquery.Selection) {
		msg = strings.Join(strings.Fields(errorhtml.Text()), " ")
		isLoggged = false
	})
	return msg, isLoggged
}

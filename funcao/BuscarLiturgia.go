package funcao

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type LeituraLiturgia struct {
	Titulo    string `json:"titulo"`
	Cor       string `json:"cor"`
	Primeira  string `json:"primeira"`
	Salmo     string `json:"salmo"`
	Segunda   string `json:"segunda"`
	Evangelho string `json:"evangelho"`
}

func trataString(texto string) string {
	tratado := strings.TrimSpace(texto)
	rx := regexp.MustCompile(`\t{1,}`)
	tratado = rx.ReplaceAllString(tratado, " ")
	return tratado
}

func BuscarLiturgia(dataMissa time.Time) LeituraLiturgia {
	// fmt.Println("Teste")
	url := fmt.Sprintf(`https://liturgiadiaria.cnbb.org.br/app/user/user/UserView.php?ano=%d&mes=%d&dia=%d`, dataMissa.Year(), dataMissa.Month(), dataMissa.Day())
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err, "Deu xabu")
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// tratamento da cor e titulo do dia liturgico

	var titulo, cor string

	cabecalho := doc.Find(".blog-header").NextFiltered(".container").Find(".bs-callout")
	titulo = trataString(cabecalho.Find("h2").Text())
	// fmt.Println(titulo)
	corObj := trataString(cabecalho.Find("p>em").Text())
	// fmt.Println(corObj)
	regexCor := regexp.MustCompile(`:\s?(\w+)$`)
	regexCorMatch := regexCor.FindAllStringSubmatch(corObj, 1)
	// fmt.Println(corObj, regexCorMatch)
	cor = regexCorMatch[0][1]

	// tratamento das leitura

	var leituras [4]string

	doc.Find("#corpo_leituras").Children().Each(func(i int, elm *goquery.Selection) {
		var sb strings.Builder
		var prefixo, sufixo string
		if i == 1 {
			// fmt.Println(elm.Find(".refrao_salmo").Text()[3:])
			prefixo = trataString(elm.Find(".title-leitura").Text())
			prefixo = prefixo[strings.Index(prefixo, "- ")+2:]
			regx := regexp.MustCompile(`Sl\s\d+([()\d]+)?`)
			prefixo = regx.FindString(prefixo)
			sufixoEl := elm.Find(".refrao_salmo")
			if sufixoEl.Text() != "" {
				sufixo = sufixoEl.Text()[3:]
			} else {
				sufixoEl = elm.Find(".REFRAO_SALMO")
				if sufixoEl.Text() != "" {
					sufixo = sufixoEl.Text()[3:]

				}

			}
			sufixo = trataString(sufixo)
		} else {
			// for {
			elm.Find("div").Find("div").Each(func(iy int, elmInterno *goquery.Selection) {
				div := elmInterno
				// fmt.Println(" -- ", div.Text())
				if valor := trataString(div.Text()); strings.Contains(valor, "Leitura") {
					prefixo = valor
				}
				if valor := trataString(div.Text()); prefixo == "" && strings.Contains(valor, "Proclama") {
					prefixo = valor
				}
				rx, err := regexp.Compile(`([\d\W]+\d\w?)+$`)
				if err != nil {
					log.Fatal("Xabu ao tratar o texto")
				}
				l := rx.FindStringIndex(prefixo)
				// fmt.Println(l)
				if prefixo != "" {
					prefixo = prefixo[0:l[0]]
				}
			})
			sufixo = trataString(elm.Find(".title-leitura").Text())
			sufixo = sufixo[strings.Index(sufixo, "- ")+2:]
		}
		sb.WriteString(prefixo)
		sb.WriteString(" - ")
		sb.WriteString(sufixo)
		leituras[i] = sb.String()

	})
	// fmt.Println(leituras)
	var obj LeituraLiturgia
	if leituras[3] != "" {
		obj = LeituraLiturgia{
			Primeira:  leituras[0],
			Salmo:     leituras[1],
			Segunda:   leituras[2],
			Evangelho: leituras[3],
			Titulo:    titulo,
			Cor:       cor,
		}
	} else {
		obj = LeituraLiturgia{
			Primeira:  leituras[0],
			Salmo:     leituras[1],
			Evangelho: leituras[2],
			Titulo:    titulo,
			Cor:       cor,
		}
	}
	fmt.Println(obj)
	return obj
}

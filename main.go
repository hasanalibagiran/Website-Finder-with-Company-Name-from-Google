package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
	"github.com/xuri/excelize/v2"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type result struct {
	Title    string
	Link     string
	Snippet  string
	Position int
}

type company struct {
	Name    string
	Country string
	Website string
}

func main() {
	//excelRead()
	//Cleaning()
}

func excelRead() {
	f, err := excelize.OpenFile("/Users/venomeux/Desktop/companies.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 4474; i < 8162; i++ {
		columnValue := strconv.Itoa(i)

		cellA := "A" + columnValue
		cellB := "B" + columnValue
		cellC := "C" + columnValue

		cellAValue, err := f.GetCellValue("Sayfa1", cellA)

		checkErr(err)

		cellBValue, err := f.GetCellValue("Sayfa1", cellB)

		checkErr(err)

		cellCValue, err := f.GetCellValue("Sayfa1", cellC)

		checkErr(err)

		companyInformation := company{Name: cellAValue, Country: cellBValue, Website: cellCValue}

		if companyInformation.Website == "" {

			googleQuery := strings.ReplaceAll(companyInformation.Name, "&", "")

			googleQuery = strings.ReplaceAll(googleQuery, " ", "+")

			googleQuery = strings.TrimSuffix(googleQuery, ".")

			googleQuery = googleQuery + "+" + strings.ToLower(companyInformation.Country)

			fmt.Println(i)
			fmt.Println(googleQuery)
			time.Sleep(3 * time.Second)
			responseFromGoogle := getData(googleQuery)

			similarityScore := 0
			similarityKing := ""
			hostnameKing := ""

			fmt.Println(len(responseFromGoogle))

			for i := 0; i < len(responseFromGoogle); i++ {
				link := responseFromGoogle[i].Link

				u, err := url.Parse(link)

				if err != nil {
					continue
				}

				rawUrl := u.Host

				splitUrlList := strings.Split(rawUrl, ".")

				hostname := ""

				if len(splitUrlList) > 2 {
					hostname = splitUrlList[1]
				} else {
					hostname = splitUrlList[0]
				}

				if hostname == "com" {
					for i := 0; i < len(splitUrlList); i++ {
						if splitUrlList[i] == "com" {
							if i >= 1 {
								hostname = splitUrlList[i-1]
							}
						}
					}
				}

				query := strings.ReplaceAll(companyInformation.Name, " ", "")

				query = strings.ToLower(query)

				similarity := similarityRatio(query, hostname)

				fmt.Println(rawUrl, hostname, similarity)

				if similarity >= 30 {
					if similarity > similarityScore {
						similarityScore = similarity
						similarityKing = rawUrl
						hostnameKing = hostname
					}
				} else if companyInformation.Country == "Türkiye" && similarity > 20 {
					if similarity > similarityScore {
						similarityScore = similarity
						similarityKing = rawUrl
						hostnameKing = hostname
					}
				}
			}
			if similarityKing != "" &&
				similarityKing != "www.instagram.com" &&
				similarityKing != "www.universidadperu.com" &&
				similarityKing != "www.abc-bahrain.com" &&
				hostnameKing != "kompass" &&
				hostnameKing != "dnb" &&
				hostnameKing != "facebook" &&
				hostnameKing != "linkedin" &&
				hostnameKing != "bloomberg" &&
				hostnameKing != "paginasamarillas" &&
				hostnameKing != "veritradecorp" &&
				hostnameKing != "wikipedia" &&
				hostnameKing != "telefoonboek" &&
				hostnameKing != "firmarehberi" &&
				hostnameKing != "goafricaonline" &&
				hostnameKing != "tradesns" &&
				hostnameKing != "opencorporates" &&
				hostnameKing != "yellowpages" &&
				hostnameKing != "piaafrica" &&
				hostnameKing != "yellow-pages" &&
				hostnameKing != "emis" &&
				hostnameKing != "datosperu" &&
				hostnameKing != "thailandbuilders" &&
				hostnameKing != "asiahoreca" &&
				hostnameKing != "bangkok-companies" {
				err = f.SetCellValue("Sayfa1", cellC, similarityKing)
				checkErr(err)
				_ = f.Save()
				fmt.Println(similarityKing)
				fmt.Println(similarityScore)
			}
			//yellow, info, company, companies, directory, go4worldbusiness, businesslist, search, data,contact, indonesiayp, israeliyp, dunsguide, iletisim, comxport, gov,5676tradeatlas
			//export 5700 5703 africanadvice, compuemprasa,directori, genealog
		}
	}

}

func similarityRatio(query string, splittedFirstLink string) int {
	return fuzzy.Ratio(query, splittedFirstLink)
}

func getData(query string) []result {
	url := "https://www.google.com/search?q=" + query

	bow := surf.NewBrowser()

	bow.SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")

	err := bow.Open(url)

	checkErr(err)

	/*request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	*/
	c := 0

	var results []result

	bow.Dom().Find("div.g").Each(func(i int, res *goquery.Selection) {
		title := res.Find("h3").First().Text()
		link, _ := res.Find("a").First().Attr("href")
		snippet := res.Find(".VwiC3b").First().Text()

		r := result{Title: title, Link: link, Snippet: snippet, Position: c + 1}

		if c < 3 {
			results = append(results, r)
		}

		c++

	})
	return results
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

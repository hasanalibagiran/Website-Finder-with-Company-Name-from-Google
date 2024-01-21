package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"net/url"
	"strconv"
	"strings"
)

func Cleaning() {
	f, err := excelize.OpenFile("/Users/venomeux/Desktop/companies.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	var wordList = []string{"instagram", "universidadperu", "abc-bahrain", "kompass", "dnb", "facebook", "linkedin", "bloomberg",
		"paginasamarillas", "veritradecorp", "wikipedia", "telefoonboek", "firmarehberi",
		"goafricaonline", "tradesns", "opencorporates", "yellowpages", "piaafrica", "yellow-pages", "emis",
		"datosperu", "thailandbuilders", "asiahoreca", "bangkok-companies", "yellow", "info", "companies", "directory", "go4worldbusiness", "businesslist", "search",
		"data", "contact", "indonesiayp", "israeliyp", "dunsguide", "iletisim", "comxport", "tradeatlas", "africanadvice", "compuempresa", "directori",
		"genealog"}

	for i := 1; i < 8163; i++ {
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
			continue
		}

		u, err := url.Parse(companyInformation.Website)

		if err != nil {
			continue
		}

		splitUrlList := strings.Split(u.String(), ".")

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

		for i := 0; i < len(wordList); i++ {
			if strings.Contains(hostname, wordList[i]) {
				err = f.SetCellValue("Sayfa1", cellC, "")
				checkErr(err)
				_ = f.Save()
				fmt.Println(hostname)
				break
			}
		}

	}
}

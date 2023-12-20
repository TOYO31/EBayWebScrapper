package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println("Erorr :", err)
	}
}
func getHtml(url string) *http.Response {
	res, err := http.Get(url)
	checkErr(err)
	if res.StatusCode > 400 {
		fmt.Println("sooemthing went  wrong...", res.StatusCode)
	}
	return res
}

func removeString(input, prefix string) string {
	if strings.HasPrefix(input, prefix) {
		return input[len(prefix):]
	}
	return input
}
func csvWriter(sData []string) {
	csvfile, err := os.OpenFile("ebaydata.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	checkErr(err)
	defer csvfile.Close()
	writer := csv.NewWriter(csvfile)
	defer writer.Flush()
	err = writer.Write(sData)
	checkErr(err)
}
func scrapData(doc *goquery.Document) {
	doc.Find("div.srp-river-results>ul.srp-results>li.s-item").Each(func(i int, item *goquery.Selection) {
		title := item.Find("div.s-item__info>a.s-item__link>div.s-item__title>span").Text()
		resTitle := removeString(title, "New Listing")
		a, _ := item.Find("a.s-item__link").Attr("href")
		price := item.Find("span.s-item__price").Text()
		scrapedData := []string{resTitle, price, a}
		csvWriter(scrapedData)

	})

}

func main() {
	for i := 1; i < 240; i++ {
		url := "https://www.ebay.com/sch/i.html?_from=R40&_trksid=p2380057.m570.l1313&_nkw=RTX&_sacat="
		newURl := url[:86] + strconv.Itoa(i)

		res := getHtml(newURl)
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		checkErr(err)
		scrapData(doc)
	}
}

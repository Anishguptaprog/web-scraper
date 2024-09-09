package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gocolly/colly"
)

type Product struct {
	Url   string
	Image string
	Name  string
	Price string
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.scrapingcourse.com", "scrapingcourse.com"),
	)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
	var products []Product
	var visitedUrls sync.Map
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		product := Product{}
		product.Url = e.ChildAttr("a", "href")
		product.Image = e.ChildAttr("img", "src")
		product.Name = e.ChildText(".product-name")
		product.Price = e.ChildText(".price")

		products = append(products, product)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})
	c.OnHTML("a", func(e *colly.HTMLElement) {
		fmt.Println(e.Attr("href"))
	})
	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scrapped!")
	})
	c.OnHTML("a.next", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		if _, found := visitedUrls.Load(nextPage); !found {
			fmt.Println("scrapping", nextPage)
			visitedUrls.Store(nextPage, struct{}{})
			e.Request.Visit(nextPage)
		}
	})

	c.OnScraped(func(r *colly.Response) {

		// open the CSV file
		file, err := os.Create("products.csv")
		if err != nil {
			log.Fatalln("Failed to create output CSV file", err)
		}
		defer file.Close()

		// initialize a file writer
		writer := csv.NewWriter(file)

		// write the CSV headers
		headers := []string{
			"Url",
			"Image",
			"Name",
			"Price",
		}
		writer.Write(headers)
		for _, product := range products {
			// convert a Product to an array of strings
			record := []string{
				product.Url,
				product.Image,
				product.Name,
				product.Price,
			}

			// add a CSV record to the output file
			writer.Write(record)
		}
		defer writer.Flush()
	})

	c.Visit("https://www.scrapingcourse.com/ecommerce")

}

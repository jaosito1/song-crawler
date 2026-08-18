package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const URL_BASE = "https://1001albumsgenerator.com"

func parseAlbumUrls(doc *html.Node) []string {
	var res []string

	var tableBody *html.Node
	for node := range doc.Descendants() {
		if node.Type == html.ElementNode && node.DataAtom == atom.Tbody {
			tableBody = node
			break
		}
	}

	var tdList []*html.Node
	for row := range tableBody.Descendants() {
		// Iterate over every tr in TableBody
		if row.Type == html.ElementNode && row.DataAtom == atom.Tr {
			// Look for the first td and save it
			for field := range row.Descendants() {
				if field.Type == html.ElementNode && field.DataAtom == atom.Td {
					tdList = append(tdList, field)
					break
				}
			}
		}
	}

	// Iterate over every td and look for the <a> element
	for _, td := range tdList {
		for inner := range td.Descendants() {
			if inner.Type == html.ElementNode && inner.DataAtom == atom.A {
				// Look in the attributes to find the href
				for _, n := range inner.Attr {
					if n.Key == "href" {
						res = append(res, n.Val)
						break
						// Exit out of the loop once href is found
					}
				}
			}
		}
	}

	return res
}

type Queue struct {
	urls []string
}

func fetchAlbumData(url string) {
	// resp, err := http.Get("URL_BASE/")
	// get item from queue
	// fetch information
	// send it back via channel
}

func main() {
	fmt.Println("Obtaining site data...")
	resp, err := http.Get(URL_BASE + "/albums")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("Parsing urls...")
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal("failed to parse html body", err)
	}

	urls := parseAlbumUrls(doc)
	if len(urls) <= 0 {
		log.Fatal("empty url slice")
	}
	fmt.Printf("Found %v album urls\n", len(urls))

	q := Queue{urls: urls}
	fmt.Println(q.urls)

	// for _, url := range q.urls {
	// 	go fetchAlbumData(url)
	// }
}

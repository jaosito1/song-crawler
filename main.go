package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const BASE_URL = "https://1001albumsgenerator.com"

func parseAlbumUrls(doc *html.Node) []string {
	var res []string

	tableBody := findByTag(doc, atom.Tbody)

	var tdList []*html.Node
	for n := range tableBody.Descendants() {
		if isTag(n, atom.Tr) {
			td := findByTag(n, atom.Td)
			tdList = append(tdList, td)
		}
	}

	for _, td := range tdList {
		a := findByTag(td, atom.A)
		href := getAttribute(a, "href")

		res = append(res, href)
	}

	return res
}


func fetchAndParse(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func fetchAlbumData(url string) {
	doc, err := fetchAndParse(BASE_URL + url)
	if err != nil {
		// ignore error or handle it (?)
	}

	node := findByClass(doc, "album-info")
	if node == nil {
		log.Fatal("missing album-info node")
	}

	data := getAlbumData(node)
	fmt.Printf("RETURNED ALBUM: %+v", data)
}

func main() {
	fmt.Println("Obtaining site data...")

	doc, err := fetchAndParse(BASE_URL + "/albums")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Parsing urls...")

	urls := parseAlbumUrls(doc)
	if len(urls) <= 0 {
		log.Fatal("empty url slice")
	}

	fmt.Printf("Found %v album urls!\n", len(urls))
}

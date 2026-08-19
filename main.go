package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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

func fetchAlbumData(urls <-chan string, data chan<- AlbumData, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range urls {
		doc, err := fetchAndParse(BASE_URL + url)
		if err != nil {
			// ignore error or handle it (?)
		}

		node := findByClass(doc, "album-info")
		if node == nil {
			log.Fatal("missing album-info node")
		}

		data <- getAlbumData(node)
	}
}

func startWorkers(urlList []string, count int) error {
	urls := make(chan string, count)
	data := make(chan AlbumData, count)
	var wg sync.WaitGroup

	file, err := NewCsvFile()
	if err != nil {
		return fmt.Errorf("failed to create csv file: %w", err)
	}

	const WORKER_COUNT = 5

	for w := 1; w <= WORKER_COUNT; w++ {
		wg.Add(1)
		go fetchAlbumData(urls, data, &wg)
	}

	for _, url := range urlList {
		urls <- url
	}
	close(urls)

	go func() {
		wg.Wait()
		close(data)
	}()

	for v := range data {
		file.WriteLine(v.Csv())
		fmt.Printf("Wrote '%v' to output.csv\n", v.Title)
	}

	return nil
}

func main() {
	fmt.Println("Obtaining site data...")

	doc, err := fetchAndParse(BASE_URL + "/albums")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parsing urls...")

	urls := parseAlbumUrls(doc)

	count := len(urls)
	if count <= 0 {
		log.Fatal("empty url slice")
	}
	fmt.Printf("Found %v album urls!\n", count)

	start := time.Now()

	fmt.Println("Retrieving album data...")
	startWorkers(urls, count)

	t := time.Now()
	elapsed := t.Sub(start)

	fmt.Printf("Process finished! Took %v seconds.", elapsed.Seconds())
}

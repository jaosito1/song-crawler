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

func worker(file *CsvFile, jobs <-chan string, counter chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range jobs {
		doc, err := fetchAndParse(BASE_URL + url)
		if err != nil {
			fmt.Printf("Error: %v", err)
			continue
		}

		node := findByClass(doc, "album-info")
		if node == nil {
			fmt.Printf("Error: %v", err)
			continue
		}

		album := getAlbumData(node)
		if err := file.WriteLine(album.Csv()); err != nil {
			fmt.Printf("Error: %v", err)
			continue
		}

		fmt.Printf("Wrote '%v' to output.csv\n", album.Title)
		counter <- 1
	}
}

func startWorkerPool(urlList []string, urlCount int, file *CsvFile) {
	var wg sync.WaitGroup

	jobs := make(chan string, urlCount)
	success := make(chan int, urlCount)

	const WORKER_LIMIT = 20

	for w := 1; w <= WORKER_LIMIT; w++ {
		wg.Add(1)
		go worker(file, jobs, success, &wg)
	}

	for _, url := range urlList {
		jobs <- url
	}
	close(jobs)

	go func() {
		// When all workers stop we stop counting succesfull tries
		wg.Wait()
		close(success)
	}()

	counter := 0
	for range success {
		counter++
	}

	fmt.Printf("Succesfully wrote %v records", counter)
}

func main() {
	fmt.Println("Obtaining site data...")

	doc, err := fetchAndParse(BASE_URL + "/albums")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parsing urls...")

	urls := parseAlbumUrls(doc)
	urlCount := len(urls)

	if urlCount <= 0 {
		log.Fatal("empty url slice")
	}

	fmt.Printf("Found %v album urls!\n", len(urls))

	file, err := NewCsvFile()
	if err != nil {
		log.Fatal("failed to create csv file: %w", err)
	}
	defer file.Close()

	start := time.Now()

	fmt.Println("Retrieving album data...")
	startWorkerPool(urls, urlCount, file)

	t := time.Now()
	elapsed := t.Sub(start)

	fmt.Printf("Process finished! Took %v seconds.", elapsed.Seconds())
}

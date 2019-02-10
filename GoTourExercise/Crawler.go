package main

import (
	"fmt"
	"sync"
)

const (
	DEPTH       = 4
	ENDINGDEPTH = DEPTH - 1
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	var visited []string
	var mutex = sync.Mutex{}
	var channels []chan int

	visited = append(visited, url)

	doCrawling(url, depth, &visited, &mutex, nil, true, &channels)

	for _, ch := range channels {
		<-ch
	}
	return
}

func Crawler(url string, depth int, visited *[]string, mutex *sync.Mutex, ch chan int) {
	doCrawling(url, depth, visited, mutex, ch, false, nil)
	checkEnd(depth, ch)
	return
}

func doCrawling(url string, depth int, visited *[]string, mutex *sync.Mutex, ch chan int, isRootCrawler bool, channels *[]chan int) {
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Found: %s %v\n", url, body)
	for _, u := range urls {
		f := getCrawlerFunction(isRootCrawler, u, depth, visited, mutex, ch, channels)
		f()
	}
}

func getCrawlerFunction(isRootCrawler bool, u string, depth int, visited *[]string, mutex *sync.Mutex, ch chan int, channels *[]chan int) func() {
	if !isRootCrawler {
		return concurrentCrawlerTask(u, depth, visited, mutex, ch)
	} else {
		return rootCrawlerTask(u, depth, visited, mutex, channels)
	}
}

func concurrentCrawlerTask(u string, depth int, visited *[]string, mutex *sync.Mutex, ch chan int) func() {
	return func() {
		if !contains(u, visited, mutex) {
			*visited = append(*visited, u)
			Crawler(u, depth-1, visited, mutex, ch)
		}
	}
}

func rootCrawlerTask(u string, depth int, visited *[]string, mutex *sync.Mutex, channels *[]chan int) func() {
	return func() {
		var ch = make(chan int)
		*channels = append(*channels, ch)
		*visited = append(*visited, u)
		go Crawler(u, depth-1, visited, mutex, ch)
	}
}

func checkEnd(depth int, ch chan int) {
	if depth == ENDINGDEPTH {
		ch <- 1
	}
}

func contains(url string, visited *[]string, mutex *sync.Mutex) bool {
	mutex.Lock()
	for _, element := range *visited {
		if element == url {
			defer mutex.Unlock()
			return true
		}
	}
	defer mutex.Unlock()
	return false
}

func main() {
	Crawl("http://golang.org/", DEPTH, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}

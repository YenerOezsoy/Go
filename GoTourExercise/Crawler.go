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

	doCrawling(url, depth, &visited, &mutex, nil, 0, &channels)

	for _, ch := range channels {
		<-ch
	}
	return
}

func Crawler(url string, depth int, visited *[]string, mutex *sync.Mutex, ch chan int) {
	doCrawling(url, depth, visited, mutex, ch, 1, nil)
	checkEnd(depth, ch)
	return
}

func doCrawling(url string, depth int, visited *[]string, mutex *sync.Mutex, ch chan int, crawlerID int, channels *[]chan int) {
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, u := range urls {
		f := getCrawlerFunction(crawlerID, u, body, depth, visited, mutex, ch, channels)
		f()
	}
}

func getCrawlerFunction(id int, u string, body string, depth int, visited *[]string, mutex *sync.Mutex, ch chan int, channels *[]chan int) func() {
	if id == 0 {
		return concurrentCrawlerSpecificTask(u, body, depth, visited, mutex, ch, channels)
	} else {
		return crawlerSpecificTask(u, body, depth, visited, mutex, ch)
	}
}

func crawlerSpecificTask(u string, body string, depth int, visited *[]string, mutex *sync.Mutex, ch chan int) func() {
	return func() {
		if !contains(u, visited, mutex) {
			mutex.Lock()
			fmt.Printf("found: %s %q\n", u, body)
			*visited = append(*visited, u)
			mutex.Unlock()
			Crawler(u, depth-1, visited, mutex, ch)
		}
	}
}

func concurrentCrawlerSpecificTask(u string, body string, depth int, visited *[]string, mutex *sync.Mutex, ch chan int, channels *[]chan int) func() {
	return func() {
		var ch = make(chan int)
		*channels = append(*channels, ch)
		fmt.Printf("found: %s \n", u)
		mutex.Lock()
		*visited = append(*visited, u)
		go Crawler(u, depth-1, visited, mutex, ch)
		mutex.Unlock()
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

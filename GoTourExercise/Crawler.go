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

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		var ch = make(chan int)
		channels = append(channels, ch)
		fmt.Printf("found: %s \n", u)
		mutex.Lock()
		visited = append(visited, u)
		go Crawler(u, depth-1, fetcher, &visited, &mutex, ch)
		mutex.Unlock()
	}

	for _, ch := range channels {
		<-ch
		close(ch)
	}
	return

}

func Crawler(url string, depth int, fetcher Fetcher, visited *[]string, mutex *sync.Mutex, ch chan int) {
	if depth <= 0 {
		checkEnd(depth, ch)
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		checkEnd(depth, ch)
		return
	}
	for _, u := range urls {
		if !contains(u, visited, mutex) {
			mutex.Lock()
			fmt.Printf("found: %s %q\n", u, body)
			*visited = append(*visited, u)
			mutex.Unlock()
			Crawler(u, depth-1, fetcher, visited, mutex, ch)
		}

	}
	checkEnd(depth, ch)

	return
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

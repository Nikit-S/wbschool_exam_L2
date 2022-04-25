package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
func getLinks(body io.Reader, host string) []string {
	var links []string
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			//todo: links list shoudn't contain duplicates
			return links
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						urlStr, err := url.Parse(attr.Val)
						if err != nil {
							continue
						}
						links = append(links, strings.TrimSpace(urlStr.Path))
					}
				}
			}

		}
	}
}

type scraper struct {
	*http.Client
	wg *sync.WaitGroup
	sync.Map
	recursive *bool
}

func (scrp *scraper) scrape(urlStr string) {

	resp, err := scrp.Get(urlStr)
	fmt.Fprintln(os.Stdout, urlStr)
	if err != nil || resp.Status != "200 OK" {
		scrp.wg.Done()
		return
	}
	var filename string
	urlStrTemp, err := url.Parse(urlStr)
	if err != nil {
		scrp.wg.Done()
		return
	}
	if urlStrTemp.Host+urlStrTemp.Path == resp.Request.URL.Host {
		os.Mkdir(resp.Request.URL.Host, os.FileMode(777))
		filename = "index.html"
	} else {

		path := resp.Request.URL.Path
		if len(path) == 1 {
			scrp.wg.Done()
			return
		}

		if path[len(path)-1] == '/' {
			filename = path[1:] + "index.html"
		} else {
			filename = path[1:]
		}
	}
	stat, err := os.Stat(filepath.Dir(resp.Request.URL.Host + "/" + filename))
	if err == nil && !stat.IsDir() {
		os.Remove(filepath.Dir(resp.Request.URL.Host + "/" + filename))
	}
	err = os.MkdirAll(filepath.Dir(resp.Request.URL.Host+"/"+filename), os.FileMode(666))
	if err != nil {
		fmt.Fprintln(os.Stderr, "os.MkdirAll error")
		fmt.Fprintln(os.Stderr, filepath.Dir(resp.Request.URL.Host+"/"+filename), err)
		scrp.wg.Done()
		return
	}
	file, err := os.Create(resp.Request.URL.Host + "/" + filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "os.Create error")
		fmt.Fprintln(os.Stderr, resp.Request.URL.Host, "/", filename, err)
		scrp.wg.Done()
		return
	}
	defer file.Close()
	if !*scrp.recursive {
		scrp.wg.Done()
		return
	}
	//fmt.Println("Creating file:", filename)
	scrp.Store(filename, struct{}{})
	r := io.TeeReader(resp.Body, file)
	for _, str := range getLinks(r, urlStrTemp.Host) {

		_, ok := scrp.Load(str)
		if ok {
			continue
		}

		scrp.Store(str, struct{}{})
		scrp.wg.Add(1)
		go scrp.scrape(urlStrTemp.Scheme + "://" + resp.Request.URL.Host + str)
	}
	scrp.wg.Done()
}

func main() {

	scrp := &scraper{}
	scrp.wg = &sync.WaitGroup{}
	scrp.Client = &http.Client{
		Timeout: 5 * time.Second,
	}

	scrp.recursive = flag.Bool("r", false, "рекурсивное скачивание")
	//reader := bufio.NewReader(os.Stdin)
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		return
	}
	scrp.wg.Add(1)
	scrp.scrape(args[0])
	//scrp.wg.Wait()
	fmt.Println(scrp.Map)
	//scrp.wg.Wait()
	scrp.wg.Wait()
	fmt.Println("Done!")

}

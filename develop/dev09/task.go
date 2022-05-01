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
func getLinks(body io.Reader) []string {
	var links []string
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
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
	fmt.Println(urlStr)
	if err != nil || resp.StatusCode != 200 {
		scrp.wg.Done()
		return
	}
	var filename string
	urlStrTemp, err := url.Parse(urlStr)
	if err != nil {
		scrp.wg.Done()
		return
	}
	// проверка на то что это корнвая папка
	if urlStrTemp.Host+urlStrTemp.Path == resp.Request.URL.Host {
		err := os.Mkdir(resp.Request.URL.Host, os.FileMode(777))
		if err != nil {
			scrp.wg.Done()
			return
		}
		filename = "index.html"
	} else {

		path := resp.Request.URL.Path
		// <a href="#">
		// <a href="/">
		if len(path) == 1 {
			scrp.wg.Done()
			return
		}

		//либо это папка и тогда мы в корень полжем index либо это файл
		if path[len(path)-1] == '/' {
			filename = path[1:] + "index.html"
		} else {
			filename = path[1:]
		}
	}

	//abc.ru + / + 1page
	//abc.ru + / + 1page/
	//abc.ru + / + //././dir1/page/
	stat, err := os.Stat(filepath.Dir(resp.Request.URL.Host + "/" + filename))

	//удаляем файл если собираемся созхдать папку с таким же назваием
	if err == nil && !stat.IsDir() {
		err = os.Remove(filepath.Dir(resp.Request.URL.Host + "/" + filename))
		if err != nil {
			scrp.wg.Done()
			return
		}
	}
	//создаем необъодимые подпапки
	err = os.MkdirAll(filepath.Dir(resp.Request.URL.Host+"/"+filename), os.FileMode(666))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "os.MkdirAll error")
		_, _ = fmt.Fprintln(os.Stderr, filepath.Dir(resp.Request.URL.Host+"/"+filename), err)
		scrp.wg.Done()
		return
	}

	//создаем файл
	file, err := os.Create(resp.Request.URL.Host + "/" + filename)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "os.Create error")
		_, _ = fmt.Fprintln(os.Stderr, resp.Request.URL.Host, "/", filename, err)
		scrp.wg.Done()
		return
	}
	defer file.Close()

	//fmt.Println("Creating file:", filename)
	scrp.Store(filename, struct{}{})
	r := io.TeeReader(resp.Body, file)
	if !*scrp.recursive {
		scrp.wg.Done()
		return
	}
	//идем посещать ссылки
	for _, str := range getLinks(r) {

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
	scrp.wg.Wait()
	fmt.Println("Done!")

}

package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"
)

func server() {

}

var client *http.Client

func loadURL(urlStr string) ([]PostData, error) {
	var URL, err = url.Parse(urlStr)
	if err != nil {
		log.Println("url parse error")
		return nil, err
	}
	parserF := selectParserF(URL)
	if parserF == nil {
		return nil, errors.New("Not found parser for " + URL.String())
	}
	resp, err := client.Get(urlStr)
	if err != nil {
		log.Println("http get error")
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("goquery parse error")
		return nil, err
	}
	doc.Url = URL
	return parserF(doc)
}

func getAllPosts() []PostData {
	const poolSize = 10
	result := make([]PostData, 0)
	urls := getURLList()
	ch := make(chan []PostData)
	wg := &sync.WaitGroup{}
	pool := make(chan struct{}, poolSize)
	done := make(chan struct{})
	go func() {
		for dat := range ch {
			result = append(result, dat...)
		}
		close(done)
	}()

	for i := range urls {
		pool <- struct{}{}
		URL := urls[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-pool }()
			log.Println("loading ", URL)
			dat, err := loadURL(URL)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			log.Printf("got %v elements from %v\n", len(dat), URL)
			ch <- dat
		}()
	}
	wg.Wait()
	close(ch)
	<-done
	return result
}

func checkDaemon() {
	posts := getAllPosts()
	log.Println("all posts ", len(posts))
	hashes := usedHashes()
	newFound := false
	newPosts := make([]PostData, 0)
	for _, post := range posts {
		hash := post.hash()
		if _, exist := hashes[hash]; exist {
			continue
		}
		newFound = true
		hashes[hash] = struct{}{}
		newPosts = append(newPosts, post)
	}
	log.Println(len(newPosts), "new posts")
	if newFound {
		saveHashes(hashes)
		sendMails(newPosts)
	}
}

func main() {

	client = &http.Client{
		Timeout: time.Second * 5,
	}

	go server()
	checkDaemon()

	return
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}

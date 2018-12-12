package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"sync"
	"time"
)

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

	inCh := make(chan ask)
	outCh := make(chan []PostData)

	//producer
	go func() {
		defer close(inCh)
		asks := copyAskList()
		for _, ask := range asks {
			inCh <- ask
		}
	}()

	wg := &sync.WaitGroup{}
	for i := 0; i < poolSize; i++ {
		//parallel worker
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ask := range inCh {
				URL := ask.GetURL()
				dat, err := loadURL(URL)
				if err != nil {
					log.Println("Error: ", err)
					return
				}

				//filter
				i := 0
				for i < len(dat) {
					delete := (dat[i].Price > ask.MaxPrice && ask.MaxPrice > 0) ||
						(dat[i].Price < ask.MinPrice && ask.MinPrice > 0)
					if delete {
						dat[i] = dat[len(dat)-1]
						dat = dat[:len(dat)-1]
					} else {
						i++
					}
				}
				outCh <- dat
			}
		}()
	}
	//monitor
	go func() {
		wg.Wait()
		close(outCh)
	}()

	//blocking consumer
	for dat := range outCh {
		result = append(result, dat...)
	}
	return result
}

func doCheck() {
	posts := getAllPosts()
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
	if newFound {
		saveHashes(hashes)
		sendMails(newPosts)
	}
}

func checkDaemon() {
	tick := time.Tick(time.Minute * 10)
	for range tick {
		doCheck()
	}
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"sync"
)

type ask struct {
	Site     string
	Part     string
	Text     string
	MaxPrice int
}

type SiteParts struct {
	Parts map[string][]string
	Names map[string]map[string]string
}

var siteParts SiteParts
var askmu sync.Mutex
var askList []ask

func initAsks() {
	askList = make([]ask, 0)
	buf, err := disk.Read("asks")
	if err != nil {
		log.Println("asks read error")
		writeAskListToDB()
		return
	}
	err = json.Unmarshal(buf, &askList)
	if err != nil {
		panic(err)
	}
}

func DeleteAsList(ind int) {
	askmu.Lock()
	defer askmu.Unlock()
	if ind < 0 || ind >= len(askList) {
		log.Println("deleted element out of range!")
		return
	}
	askList = append(askList[:ind], askList[ind+1:]...)
	go writeAskListToDB()
}

func AppendAskList(a ask) {
	askmu.Lock()
	defer askmu.Unlock()
	askList = append(askList, a)
	go writeAskListToDB()
}

func writeAskListToDB() {
	buf, err := json.Marshal(askList)
	if err != nil {
		return
	}
	disk.Write("asks", buf)
}

func NewAsk(site, part, text string, maxPrice int) (ask, error) {
	if !validSitePart(site, part) {
		return ask{}, errors.New("Non valid site-part pair ")
	}
	return ask{
		Site:     site,
		Part:     part,
		Text:     text,
		MaxPrice: maxPrice,
	}, nil
}

func validSitePart(site, part string) bool {
	if site != "avito" && site != "youla" {
		log.Println("wrong site ", site)
		return false
	}
	if names, exist := siteParts.Names[site]; !exist {
		return false
	} else {
		if _, exist := names[part]; !exist {
			return false
		}
	}
	return true
}

func copyAskList() []ask {
	askmu.Lock()
	defer askmu.Unlock()
	result := make([]ask, len(askList))
	copy(result, askList)
	return result
}

func (a ask) GetURL() string {
	switch a.Site {
	case "avito":
		return avitoURL(a)
	case "youla":
		return youlaURL(a)
	default:
		return ""
	}
}

func avitoURL(a ask) string {
	text := url.QueryEscape(a.Text)
	return fmt.Sprintf("http://www.avito.ru/moskva/%v?q=%v", a.Part, text)
}

func youlaURL(a ask) string {
	text := url.QueryEscape(a.Text)
	return fmt.Sprintf("http://youla.ru/moskva/%v?q=%v", a.Part, text)
}

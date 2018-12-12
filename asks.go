package main

import (
	"encoding/json"
	"errors"
	"log"
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
	if ind < 0 || ind >= len(askList) {
		log.Println("deleted element out of range!")
		return
	}
	askList = append(askList[:ind], askList[ind+1:]...)
	go writeAskListToDB()
}

func AppendAskList(a ask) {
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

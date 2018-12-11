package main

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/peterbourgon/diskv"
	"net/http"
	"strings"
)

type SiteParts struct {
	Parts map[string][]string
	Names map[string]map[string]string
}

var parts SiteParts

func add(site, href, name string) {
	parts.Parts[site] = append(parts.Parts[site], href)
	names := parts.Names[site]
	if names == nil {
		names = make(map[string]string)
	}
	names[href] = name
	parts.Names[site] = names
}

func main() {
	parts.Parts = make(map[string][]string)
	parts.Names = make(map[string]map[string]string)

	parseAvitoParts()
	parseYoulaParts()
	save()
}

func parseAvitoParts() {
	buf, err := http.Get("http://www.avito.ru/moskva/")
	if err != nil {
		panic(err)
	}
	defer buf.Body.Close()

	doc, err := goquery.NewDocumentFromReader(buf.Body)
	if err != nil {
		panic(err)
	}
	doc.Find("a.category-map-link").Each(func(i int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		if !exist {
			return
		}
		ind := strings.LastIndex(href, "/")
		if ind > -1 {
			href = href[ind+1:]
		}
		ind = strings.Index(href, "?")
		if ind > -1 {
			href = href[:ind]
		}
		name := strings.TrimSpace(s.Text())
		add("avito", href, name)
	})
}

func parseYoulaParts() {
	buf, err := http.Get("http://www.youla.ru/moskva/")
	if err != nil {
		panic(err)
	}
	defer buf.Body.Close()

	doc, err := goquery.NewDocumentFromReader(buf.Body)
	if err != nil {
		panic(err)
	}
	_ = doc
	//We can't parse parts from youla, cz JS, sojust skip it
	return
	//
}

func save() {
	disk := diskv.New(diskv.Options{BasePath: "."})
	buf, err := json.Marshal(parts)
	if err != nil {
		panic(err)
	}
	disk.Write("parts", buf)
}

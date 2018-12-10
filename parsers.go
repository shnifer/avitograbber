package main

import (
	"crypto/md5"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type PostData struct {
	Title string
	Price int
	Href  string
}

func (p PostData) hash() [md5.Size]byte {
	return md5.Sum(append([]byte(p.Title+p.Href), byte(p.Price)))
}

type ParserF = func(*goquery.Document) ([]PostData, error)

func YoulaParser(doc *goquery.Document) (result []PostData, err error) {
	BaseURL := doc.Url
	doc.Find("li.product_item").Each(func(ind int, s *goquery.Selection) {
		a := s.Find("a").First()
		div := s.Find("div.product_item__description ").First()
		if a.Size() == 0 || div.Size() == 0 {
			log.Println("a and div not found")
			return
		}
		if div.Nodes[0].FirstChild == nil {
			log.Println("div.Nodes[0].FirstChild==nil")
			return
		}
		title, exist := a.Attr("title")
		if !exist {
			log.Println("no title")
			return
		}
		href, exist := a.Attr("href")
		if !exist {
			log.Println("no href")
			return
		}
		linkUrl, err := BaseURL.Parse(href)
		if err != nil {
			log.Println("can't parse href ", href)
			return
		}
		href = linkUrl.String()
		priceStr := div.Nodes[0].FirstChild.Data
		priceStr = strings.Replace(strings.TrimSpace(priceStr), " ", "", -1)
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			log.Println("strange price ", priceStr)
			return
		}
		result = append(result, PostData{
			Title: title,
			Price: price,
			Href:  href,
		})
	})
	return result, nil
}

func AvitoParser(doc *goquery.Document) (result []PostData, err error) {
	BaseURL := doc.Url
	doc.Find("div.item_table-header").Each(func(ind int, s *goquery.Selection) {
		a := s.Find("a.item-description-title-link").First()
		span := s.Find("span.price").First()
		if a.Size() == 0 || span.Size() == 0 {
			log.Println("a and span not found")
			return
		}
		title := a.Text()
		href, exist := a.Attr("href")
		if !exist {
			log.Println("no href")
			return
		}
		linkUrl, err := BaseURL.Parse(href)
		if err != nil {
			log.Println("can't parse href ", href)
			return
		}
		href = linkUrl.String()
		priceStr, exist := span.Attr("content")
		if !exist {
			log.Println("no price content")
			return
		}
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			log.Println("strange price ", priceStr)
		}

		result = append(result, PostData{
			Title: title,
			Price: price,
			Href:  href,
		})
	})
	return result, nil
}

func selectParserF(url *url.URL) ParserF {
	host := url.Host
	switch {
	case strings.Contains(host, "avito.ru"):
		return AvitoParser
	case strings.Contains(host, "youla.ru"):
		return YoulaParser
	default:
		return nil
	}
}

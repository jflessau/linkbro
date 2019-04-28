package main

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)

func checkLoop(list []Link, url string, all bool, n int, max int, lat int) ([]Link, string, bool, int, int, int) {
	showProgress(list)
	time.Sleep(time.Duration(lat) * time.Millisecond)

	// get link from list to check
	var l = Link{CheckStatus: false, CheckPage: false}
	for _, entry := range list {
		if entry.CheckStatus || entry.CheckPage {
			l = entry
			break
		}
	}

	// end if no links to check are left
	if !l.CheckStatus {
		return list, url, all, n, max, lat
	}

	// catch mailto links
	if isMailTo(l.Href) {
		list = setLinkProps(l, list, 0, false, false)
		return checkLoop(list, url, all, n, max, lat)
	}

	// fetch page
	res, err := http.Get(l.Href)
	if err != nil {
		log.Warning(err)
		list = setLinkProps(l, list, 0, false, false)
		return checkLoop(list, url, all, n, max, lat)
	}
	StatusCode := res.StatusCode
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Warning("reading DOM failed")
		list = setLinkProps(l, list, StatusCode, false, false)
		return checkLoop(list, url, all, n, max, lat)
	}
	defer res.Body.Close()

	// catch non 200 status codes
	if StatusCode != 200 {
		log.Notice("non 200 at", l.Href)
		list = setLinkProps(l, list, StatusCode, false, false)
		return checkLoop(list, url, all, n, max, lat)
	}

	// start over with next link if current one is not on the specified domain
	if !l.CheckPage {
		list = setLinkProps(l, list, 200, false, false)
		return checkLoop(list, url, all, n, max, lat)
	}

	// mark link as fully checked
	list = setLinkProps(l, list, StatusCode, false, false)

	// add links on page to list
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok && !isOnList(prepareHref(href, l.Href, url), list) && (all || url == l.Href) && n < (max-1) {
			n++
			list = append(list, Link{
				Href:        prepareHref(href, l.Href, url),
				CheckStatus: true,
				CheckPage:   sameDomain(prepareHref(href, l.Href, url), url),
				StatusCode:  0,
				Origin:      l.Href})
		}
	})

	return checkLoop(list, url, all, n, max, lat)
}

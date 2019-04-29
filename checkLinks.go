package main

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)

func checkLoop(list []Link, args LoopCheckArgs) ([]Link, LoopCheckArgs) {
	showProgress(list)
	time.Sleep(time.Duration(args.Lat) * time.Millisecond)

	// get link from list to check
	var l = Link{CheckStatus: false, CheckPage: false}
	for _, entry := range list {
		if entry.CheckStatus || entry.CheckPage {
			l = entry
			break
		}
	}

	// end if no links to check are left
	if !l.CheckStatus { return list, args }

	// catch mailto links
	if isMailTo(l.Href) {
		list = setLinkProps(l, list, 0, false, false)
		return checkLoop(list, args)
	}

	// request page
	res, err := http.Get(l.Href)
	if err != nil {
		log.Warning(err)
		list = setLinkProps(l, list, 0, false, false)
		return checkLoop(list, args)
	}
	StatusCode := res.StatusCode
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Warning("reading DOM failed")
		list = setLinkProps(l, list, StatusCode, false, false)
		return checkLoop(list, args)
	}
	defer res.Body.Close()

	// catch non 200 status codes
	if StatusCode != 200 {
		log.Notice("non 200 at", l.Href)
		list = setLinkProps(l, list, StatusCode, false, false)
		return checkLoop(list, args)
	}

	// mark link as checked
	list = setLinkProps(l, list, StatusCode, false, false)

	// switch to next link if current one is no internal link
	if !l.CheckPage { return checkLoop(list, args) }

	// add links on page to list
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok && !isOnList(prepareHref(href, l.Href, args.Url), list) && (args.All || args.Url == l.Href) && args.Num < (args.Max-1) {
			args.Num++
			list = append(list, Link{
				Href:        prepareHref(href, l.Href, args.Url),
				CheckStatus: true,
				CheckPage:   sameDomain(prepareHref(href, l.Href, args.Url), args.Url),
				StatusCode:  0,
				Origin:      l.Href})
		}
	})

	return checkLoop(list, args)
}

package main

import (
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"net/http"
	"strings"
)

func (r Redacted) Redacted() interface{} {
	return logging.Redact(string(r))
}

func checkInput(args LoopCheckArgs) (LoopCheckArgs, error) {
	if args.Max == 0 {
		return args, errors.New("please set n > 1")
	}
	if args.Url == "" {
		return args, errors.New("Please specify an url by using the -url flag")
	}
	if !testUrl(args.Url) {
		return args, errors.New("Fetch failed. Don’t forget to specify the protocol: http:// or https://")
	}

	log.Notice("------------")
	if args.Max > 1000 {
		oldN := args.Max
		args.Max = 1000
		log.Warningf("Max Link amount was set to %d. %d is a bit too much, huh?", args.Max, oldN)
	}
	if args.Lat < 250 {
		oldLat := args.Lat
		args.Lat = 250
		log.Warningf("Latency was set to %dms. %dms is a bit dangerous, huh?", args.Lat, oldLat)
	}

	return args, nil
}

func setLinkProps(l Link, list []Link, statusCode int, checkStatus, checkPage bool) []Link {
	for n := 0; n < len(list); n++ {
		if list[n].Href == l.Href {
			list[n].CheckStatus = checkStatus
			list[n].CheckPage = checkPage
			list[n].StatusCode = statusCode
			return list
		}
	}
	log.Panic("list is broken")
	return list
}

func sameDomain(urlOne, urlTwo string) bool {
	return strings.ToLower(getDomain(urlOne, false)) == strings.ToLower(getDomain(urlTwo, false))
}

func isOnList(href string, list []Link) bool {
	for _, entry := range list {
		if entry.Href == href {
			return true
		}
	}
	return false
}

func getDomain(url string, protocol bool) string {
	protocolStr := ""
	if protocol {
		protocolStr = "http://"
		if strings.Index(url, "https:") >= 0 {
			protocolStr = "https://"
		}
	}

	url = strings.Replace(url, "http://", "", -1)
	url = strings.Replace(url, "https://", "", -1)

	var offset, end int
	if strings.Index(url, "//") >= 0 {
		offset = 2
	}
	urlCut := strings.Replace(url, "//", "", -1)
	end = strings.Index(urlCut, "/")
	if end >= 0 {
		return protocolStr + string(url[0:(end+offset)])
	}

	return protocolStr + url
}

func prepareHref(href, originHref, url string) string {
	if sameDomain(href, url) {
		return href
	}
	if len(href) < 1 {
		return originHref
	}
	if string(href[0]) == "/" {
		return getDomain(url, true) + href
	}
	if string(href[0]) == "#" {
		return originHref + addSlash(originHref) + href
	}
	return href
}

func addSlash(s string) string {
	if string(s[len(s)-1]) == "/" {
		return ""
	}
	return "/"
}

func isMailTo(url string) bool {
	return strings.Index(url, "mailto:") == 0
}

func showProgress(list []Link) {
	var n = len(list)
	var done = 0
	for _, l := range list {
		if !l.CheckStatus && !l.CheckPage {
			done++
		}
	}
	log.Info((fmt.Sprintf("\r\r(%d of %d links checked)", done, n)))
}

func printResult(list []Link, args LoopCheckArgs) {
	var ok, notSure, notOk int
	for n, l := range list {
		switch l.StatusCode {
		case 200:
			ok++
			if args.Slim {
				break
			}
			log.Infof("%d: ok: %s ---> %s", (n + 1), l.Origin, l.Href)
		case 0:
			notSure++
			log.Noticef("%d: not sure: %s ---> %s", (n + 1), l.Origin, l.Href)
		default:
			notOk++
			log.Warningf("%d: not ok: %s ---> %s", (n + 1), l.Origin, l.Href)
		}
	}

	log.Notice("\n------------\nSUMMARY")
	log.Infof("ok:       %d", ok)
	log.Noticef("not sure: %d", notSure)
	log.Warningf("not ok:   %d", notOk)
	log.Notice("------------\n")
}

func testUrl(url string) bool {
	_, err := http.Get(url)
	return err == nil
}

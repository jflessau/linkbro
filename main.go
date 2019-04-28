package main

import (
	"os"
  "flag"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main logger")
var format = logging.MustStringFormatter(
	`%{color}%{message}%{color:reset} `,
)

func main() {
  file, _ := os.OpenFile("./log/log.log", os.O_RDWR|os.O_CREATE, 0666)
  backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFile := logging.NewLogBackend(file, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFile, backendFormatter)

  all := flag.Bool("a", false, "scan entire domain")
  url := flag.String("url", "", "domain to perform a check on")
  n := flag.Int("n", 50, "amount of links to check")
  lat := flag.Int("l", 1000, "time [ms] to wait between fetches")
  flag.Parse()

  if *n == 0 {
    log.Warning("please set n > 1")
    return
  }

  if *url == "" {
    log.Warning("please specify an url with the -url flag")
    return
  }

	log.Notice("---")

	if *n > 300 {
		oldN := *n
		n = new(int)
		*n = 300
		log.Warningf("n was set to %d, %d exceeds the limit", *n, oldN)
	}

	if *lat < 250 {
		oldLat := *n
		lat = new(int)
		*lat = 250
		log.Warningf("latency was set to %dms, %dms is a bit dangerous, huh?", *lat, oldLat)
	}

	log.Noticef("url:          %s", *url)
	log.Noticef("depth:        %d links", *n)
	log.Noticef("whole domain: %t", *all)
	log.Noticef("latency:      %d milliseconds", *lat)

	log.Warningf("---\n\n")

	if !testUrl(*url) {
		log.Warning("Fetch failed. Don’t forget to specify the protocol: http:// or https://")
		return
	}


	list := []Link{Link{Href: *url, Origin: *url, CheckPage: true, CheckStatus: true}}
  list, _, _, _, _, _ = checkLoop(list, *url, *all, 0, *n, *lat)
	printResult(list)
}

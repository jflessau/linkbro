package main

import (
	"flag"
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("")
var format = logging.MustStringFormatter(`%{color}%{message}%{color:reset} `)

func main() {
	os.Remove("./log/log.log")
	file, _ := os.OpenFile("./log/log.log", os.O_RDWR|os.O_CREATE, 0666)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFile := logging.NewLogBackend(file, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFile, backendFormatter)
	all := flag.Bool("a", false, "scan entire domain")
	slim := flag.Bool("s", false, "scan entire domain")
	url := flag.String("url", "", "domain to perform a check on")
	n := flag.Int("n", 50, "amount of links to check")
	lat := flag.Int("l", 1000, "time [ms] to wait between fetches")
	flag.Parse()

	args := LoopCheckArgs{
		Url:  *url,
		All:  *all,
		Num:  0,
		Max:  *n,
		Lat:  *lat,
		Slim: *slim}

	var err error
	args, err = checkInput(args)
	if err != nil {
		log.Error(err)
		return
	}

	log.Noticef("url:          %s", args.Url)
	log.Noticef("depth:        %d links", args.Max)
	log.Noticef("whole domain: %t", args.All)
	log.Noticef("latency:      %d milliseconds", args.Lat)
	log.Noticef("------------\n")

	// perform check
	list := []Link{Link{Href: *url, Origin: *url, CheckPage: true, CheckStatus: true}}
	list, _ = checkLoop(list, args)
	printResult(list, args)
}

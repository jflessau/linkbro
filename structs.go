package main

type Redacted string

type Link struct {
	Href        string
	CheckStatus bool
	CheckPage   bool
	StatusCode  int
	Origin      string
}

type LoopCheckArgs struct {
	Url  string
	All  bool
	Num  int
	Max  int
	Lat  int
	Slim bool
}

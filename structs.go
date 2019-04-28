package main

type Redacted string

type Link struct {
  Href        string
  CheckStatus bool
  CheckPage   bool
  StatusCode  int
  Origin      string
}

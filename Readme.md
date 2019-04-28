
# linkbro
Check all public links on a site and discover broken links.

## Build & Run
1. Install go packages with `$ go get ./...`
2. Use `$ go build -o linkbro` to build a binary
3. Make it executable with `$ chmod +x ./linkbro`
4. Run with `$ ./linkbro`

## Usage
```
$ linkbro -url https://example.com -n 42 -l 500 -a -s
```

### Options
| Option | Type    | Optional | Default | Description                           |
|:------ |:------- | -------- | ------- | ------------------------------------- |
| -url   | string  | no       | -       | url to start checks on                |
| -n     | integer | yes      | 50      | max number of links to check          |
| -l     | integer | yes      | 1000    | latency between requests              |
| -a     | boolean | yes      | false   | check all links on a domain             |
| -s     | boolean | yes      | false   | slim summary (show broken links only) |


## Hints
Without the `-a` flag, only links on the specified page are checked. For example:
```
$ linkbro -url https://example.com
```
checks internal and external links found on `https://example.com`.
```
$ linkbro -url https://example.com -a
```
will also check all links discovered by hopping from link to link within the domain `example.com`.


## Example
`linkbro -url example.com -n 10` will deliver: the following output:
```
------------
url:          https://example.com
depth:        10 links
whole domain: false
latency:      1000 milliseconds
------------

(0 of 1 links checked)
(1 of 2 links checked)
(2 of 2 links checked)
1: ok: https://example.com ---> https://example.com
2: ok: https://example.com ---> http://www.iana.org/domains/example

------------
SUMMARY
ok:       2
not sure: 0
not ok:   0
------------
```

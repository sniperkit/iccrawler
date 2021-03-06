package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"

	"github.com/stanxii/iccrawler/crawlerPassive/config"
	"github.com/stanxii/iccrawler/crawlerPassive/links"
)

func worker(l *links.Links, seeds []string) {

	//real task worker very long time need be cancel up code.
	fmt.Println("Looping...working............... do work worker## ", time.Now())
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel() //cancel when we are finished conxuming string list.!

	// defer close(Pages)

	// ctx, cancel := context.WithCancel(context.Background())

	List := l.ListURLS(ctx, seeds)

	Pages := l.DetailURLS(ctx, List)

	Storages := l.DetailPage(ctx, Pages)

	l.StorageCSV(ctx, Storages)
	//wait for pages close then close storage channel_price

	elapsed := time.Since(start)
	fmt.Println("Looping...working End End......... All storaged finish consumming save to db...  It took: ", elapsed)
}

func main() {

	fmt.Println("hello ..")

	configFile := flag.String("config", "./etc/crawler.passive.toml", "crawler passive config file.")
	flag.Parse()

	fmt.Printf("configFile =: %v\n", *configFile)

	c, err := config.NewConfigWithFile(*configFile)
	if err != nil {
		log.Fatal("Err config file error.", err)
	}

	l := links.NewLinks(c)
	if err != nil {
		log.Fatal("Err crawler Passive error.", err)
	}

	defer l.Stop()

	//passive seeds
	seeds := []string{
		"http://www.szlcsc.com/so/catalog/308.html",
		"http://www.szlcsc.com/so/catalog/312.html",
		"http://www.szlcsc.com/so/catalog/316.html",
		"http://www.szlcsc.com/so/catalog/348.html",
	}

	fmt.Println("seeds len=", len(seeds))
	//first time
	worker(l, seeds)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Println()
	fmt.Println(sig)
	l.Stop()
	fmt.Println("awaiting signal")

	//send to nsq exist for loop.
	fmt.Println("exiting")
}

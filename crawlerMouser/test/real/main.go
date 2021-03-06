package main

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/stanxii/iccrawler/crawlerSZLC/links"
)

func worker(ctx context.Context, l *links.Links, seeds []string, out chan<- string) error {

	//real task worker very long time need be cancel up code.
	fmt.Println("### Default do work worker## ", time.Now())

	err := l.CrawlerSZLC(ctx, seeds, out)
	defer close(out)

	if err != nil {
		fmt.Println("dowork err :", err)
	}

	///////////block worker until finish or canceled.
	select {
	case <-ctx.Done():
		fmt.Println("one day finished.")
		return nil
	}

}

func main() {

	//ten years.
	stop := time.After(3650 * 24 * time.Second)

	dur := 24
	//tick for one day run once worker
	tick := time.NewTicker(time.Duration(dur) * time.Second)
	defer tick.Stop()

	////////////////////////////////start/////////////
	l := links.NewLinks()
	defer l.Stop()
	u := "http://127.0.0.1:8001/szlcsccat?keyword=2222"
	// u := "http://10.8.15.9:8001/szlcsccat?keyword=2222"
	seeds, err := l.GetSeedURLS(u)
	if err != nil {
		fmt.Println("get init seed error")
		log.Fatal(err)
	}
	fmt.Println("seeds len=", len(seeds))
	// done := make(chan struct{})

	for {

		ctx2, cancel2 := context.WithCancel(context.Background())

		List := make(chan string)
		Pages := make(chan string)
		Storages := make(chan interface{})

		go l.DetailURLS(ctx2, Pages, List)
		go l.DetailPage(ctx2, Storages, Pages)
		go l.StorageCockDB(ctx2, Storages)

		durS := 2*dur - 4
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(durS)*time.Second)
		select {
		case <-stop:
			fmt.Println("all done.")
			//close
			cancel2()

			//wait for storage finish
			close(Storages)

			cancel()
			return
		case <-tick.C:
			go worker(ctx, l, seeds, List)
		}
	}

}

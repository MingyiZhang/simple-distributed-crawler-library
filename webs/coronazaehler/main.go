package main

import (
	"flag"
	"strings"

	"distributed-crawler-demo/config"
	"distributed-crawler-demo/duplicate/client"
	"distributed-crawler-demo/engine"
	saver "distributed-crawler-demo/persist/client"
	"distributed-crawler-demo/rpchelper"
	"distributed-crawler-demo/scheduler"
	"distributed-crawler-demo/webs/coronazaehler/parser"
	worker "distributed-crawler-demo/worker/client"
)

var (
	itemSaverHost = flag.String("itemsaver_host", "", "itemsaver host")
	workerHosts   = flag.String("worker_hosts", "", "worker hosts (comma separated)")
	checkerHost   = flag.String("checker_host", "", "duplicate checker host")
)

func main() {
	flag.Parse()
	itemChan, err := saver.ItemSaver(*itemSaverHost)
	if err != nil {
		panic(err)
	}

	pool := rpchelper.CreateClientPool(strings.Split(*workerHosts, ","))
	processor := worker.CreateProcessor(pool)

	checkerClient, err := rpchelper.NewClient(*checkerHost)
	checker := client.CreateChecker(checkerClient)

	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      1,
		ItemChan:         itemChan,
		RequestProcessor: processor,
		DuplicateChecker: checker,
	}

	e.Run(engine.Request{
		Url:    "https://www.coronazaehler.de",
		Parser: engine.NewFuncParser(parser.ParseCounty, config.ParseCounty),
	})
}

package main

import (
  "flag"
  "strings"

  "distributed-crawler-demo/config"
  "distributed-crawler-demo/engine"
  saver "distributed-crawler-demo/persist/client"
  "distributed-crawler-demo/rpchelper"
  "distributed-crawler-demo/scheduler"
  "distributed-crawler-demo/webs/mockweb/parser"
  worker "distributed-crawler-demo/worker/client"
)

const cityListUrl = "http://localhost:8080/mock/www.zhenai.com/zhenghun"

var (
  itemSaverHost = flag.String("itemsaver_host", "", "itemsaver host")
  workerHosts = flag.String("worker_hosts", "", "worker hosts (comma separated)")
)

func main() {
  flag.Parse()
  itemChan, err := saver.ItemSaver(*itemSaverHost)
  if err != nil {
    panic(err)
  }

  pool := rpchelper.CreateClientPool(strings.Split(*workerHosts, ","))
  processor := worker.CreateProcessor(pool)

  e := engine.ConcurrentEngine{
    Scheduler:        &scheduler.QueuedScheduler{},
    WorkerCount:      100,
    ItemChan:         itemChan,
    RequestProcessor: processor,
  }

  e.Run(engine.Request{
    Url:    cityListUrl,
    Parser: engine.NewFuncParser(parser.ParseCityList, config.ParseCityList),
  })
}

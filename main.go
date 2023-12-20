package main

import (
	"context"
	"flag"
	"github.com/ponishaofficial/darkoob/utils"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var scenarioFolder string

func init() {
	flag.StringVar(&scenarioFolder, "scenarios", "./scenarios", "location of scenarios folder")
	flag.Parse()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go handleInterrupts(ctx, cancel)
	yamls, err := utils.ReadYAMLFiles(scenarioFolder)
	if err != nil {
		log.Fatal(err)
	}

	wg := new(sync.WaitGroup)
	statCh := make(chan *utils.Stats)

	go utils.DoStats(statCh)

	for file := range yamls {
		if yamls[file].Name == "" {
			yamls[file].Name = file
		}
		if yamls[file].Iteration < 0 {
			yamls[file].Iteration = 0
		}
		if yamls[file].Concurrency < 1 {
			yamls[file].Concurrency = 1
		}

		for i := 1; i <= yamls[file].Iteration; i++ {
			wg.Add(1)
			go utils.RunScenario(ctx, wg, statCh, i, yamls[file])
		}
	}

	wg.Wait()
	close(statCh)
	utils.ShowStats(yamls)
}

func handleInterrupts(ctx context.Context, cancel context.CancelFunc) {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	select {
	case <-ctx.Done():
		return
	case <-sig:
		cancel()
		return
	}
}

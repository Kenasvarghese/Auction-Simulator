package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	auctionhouse "github.com/Kenasvarghese/Auction-Simulator/auction_house"
	"github.com/Kenasvarghese/Auction-Simulator/bidders"
	"github.com/Kenasvarghese/Auction-Simulator/config"
)

func main() {
	cfg := config.LoadConfig()
	err := cfg.Validate()
	if err != nil {
		log.Fatal(fmt.Errorf("invalid configuration: %w", err))
		return
	}
	fmt.Printf("Config: %+v\n", cfg)
	// set the max procs to the configured vcpu
	runtime.GOMAXPROCS(cfg.VCPU)
	// create the bidder list
	bidderList := make([]bidders.Bidder, 0, cfg.NumBidders)
	for i := range cfg.NumBidders {
		bidder := bidders.NewBidder(i)
		bidderList = append(bidderList, bidder)
	}
	ah := auctionhouse.NewAuctionHouse(bidderList)
	var wt sync.WaitGroup

	start := time.Now()
	log.Printf("Starting %d auctions with %d bidders\n", cfg.NumAuctions, cfg.NumBidders)

	// create semaphores for VCPU and Memory
	cpuSem := make(chan struct{}, cfg.VCPU)   // semaphore to limit concurrent auctions to VCPU count
	memSem := make(chan struct{}, cfg.Memory) // semaphore to limit memory usage (1 token = 1 MB)

	// assume each auction uses 10 MB memory and 1 VCPU

	// launch the auctions as goroutines
	for i := range cfg.NumAuctions {
		for range cfg.AuctionVCPU {
			cpuSem <- struct{}{} // acquire a VCPU token
		}

		// acquire 10 MB tokens before starting auction
		for range cfg.AuctionMemory {
			memSem <- struct{}{}
		}
		wt.Go(func() {
			ah.RunAuction(i+1, cfg.AuctionTimeoutMs)
			<-cpuSem // release VCPU token
			// release memory tokens after auction
			for range cfg.AuctionMemory {
				<-memSem
			}
		})
	}
	wt.Wait()
	end := time.Now()
	log.Printf("Completed %d auctions in %v\n", cfg.NumAuctions, end.Sub(start))
}

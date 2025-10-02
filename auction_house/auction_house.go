package auctionhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Kenasvarghese/Auction-Simulator/bidders"
	"github.com/Kenasvarghese/Auction-Simulator/domain"
	"github.com/Kenasvarghese/Auction-Simulator/utils"
)

type auctionHouse struct {
	bidderList []bidders.Bidder
}

func NewAuctionHouse(bidderList []bidders.Bidder) AuctionHouse {
	return &auctionHouse{
		bidderList: bidderList,
	}
}

// runAuction runs a single auction with the given id and timeout (ms). It writes an output file and returns result.
func (a *auctionHouse) RunAuction(auctionID int, timeoutMs int) domain.AuctionResult {
	attrs := utils.MakeAttributes()
	start := time.Now()

	log.Printf("Auction %d started with timeout %d ms and %d bidders", auctionID, timeoutMs, len(a.bidderList))

	// create a per-auction context representing the timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()

	ch := make(chan domain.Bid, len(a.bidderList)) // buffered to avoid blocked sends
	var wg sync.WaitGroup

	// Fan-out to bidders
	for i := range a.bidderList {
		// We'll launch bidders as goroutines
		wg.Go(func() {
			// simulate bidder behavior
			bidder := a.bidderList[i]
			bidder.Call(ctx, attrs, ch)
		})
	}

	// collector goroutine: waits for all bidder goroutines to finish OR for context timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// determine winner (highest price)
	var winner domain.Bid

	bids := make([]domain.Bid, 0, len(a.bidderList))
loop:
	for {
		select {
		case b := <-ch:
			bids = append(bids, b)
			if winner.Price < b.Price {
				winner = b
			}

		case <-done:
			// all bidders finished
			break loop
		case <-ctx.Done():
			// timeout reached: drain any immediate ready bids from ch (non-blocking)
		drain:
			for {
				select {
				case b := <-ch:
					bids = append(bids, b)
					if winner.Price < b.Price {
						winner = b
					}
				default:
					break drain
				}
			}
			break loop
		}
	}

	end := time.Now()
	// calculate auction duration in ms
	durationMs := end.Sub(start).Milliseconds()

	result := domain.AuctionResult{
		AuctionID:       auctionID,
		StartTime:       start,
		EndTime:         end,
		DurationMs:      durationMs,
		NumBidders:      len(a.bidderList),
		NumBidsReceived: len(bids),
		Bids:            bids,
		Winner:          winner,
		Attributes:      attrs,
	}

	// write output JSON file
	if err := os.MkdirAll("./Output", 0755); err == nil {
		fn := fmt.Sprintf("%s/auction_%03d.json", "./Output", auctionID)
		f, err := os.Create(fn)
		if err == nil {
			enc := json.NewEncoder(f)
			enc.SetIndent("", "  ")
			_ = enc.Encode(result)
			f.Close()
		} else {
			log.Printf("failed to create output %s: %v", fn, err)
		}
	} else {
		log.Printf("failed to create outputs dir: %v", err)
	}
	return result
}

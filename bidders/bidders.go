package bidders

import (
	"context"
	"math/rand"
	"time"

	"github.com/Kenasvarghese/Auction-Simulator/domain"
)

type bidder struct {
	bidderID int
}

func NewBidder(i int) Bidder {
	return &bidder{
		bidderID: i,
	}
}

// Call simulates a call to the bidder
func (b *bidder) Call(ctx context.Context, attrs map[string]string, ch chan<- domain.Bid) {
	// here we are simulating bidder behavior
	// in a real system this would be an HTTP call or similar
	simulateBid(ctx, b.bidderID, attrs, ch)
}

// simulateBidder simulates a bidder that may or may not respond within the auction context.
// It sends a Bid to ch if it responds before ctx is done.
func simulateBid(ctx context.Context, bidderID int, attrs map[string]string, ch chan<- domain.Bid) {
	// for simplicity we ignore attributes for the bidder
	_ = attrs

	// Randomly decide whether this bidder will respond (70% chance respond, adjust as needed)
	respond := rand.Intn(100) < 70
	// Random network delay between 10 and 200 ms
	delay := time.Duration(rand.Intn(191)+10) * time.Millisecond

	select {
	case <-time.After(delay):
		if !respond {
			return // no response
		}
		// Create a random bid price between 0.01 and 5.00 (CPM)
		price := rand.Float64()*4.99 + 0.01
		bid := domain.Bid{
			BidderID:  bidderID,
			Price:     price,
			LatencyMs: int(delay / time.Millisecond),
		}
		select {
		case ch <- bid:
		case <-ctx.Done():
			// auction closed before we could send
			return
		}
	case <-ctx.Done():
		// auction timed out before bidder responded
		return
	}
}

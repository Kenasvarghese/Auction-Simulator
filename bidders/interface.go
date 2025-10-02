package bidders

import (
	"context"

	"github.com/Kenasvarghese/Auction-Simulator/domain"
)

// Bidder defines the interface for a bidder in the auction simulation.
type Bidder interface {
	Call(ctx context.Context, attrs map[string]string, ch chan<- domain.Bid)
}

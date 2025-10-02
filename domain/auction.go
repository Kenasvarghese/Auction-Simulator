package domain

import "time"

// AuctionResult represents the result of an auction
type AuctionResult struct {
	AuctionID       int               `json:"auction_id"`
	StartTime       time.Time         `json:"start_time"`
	EndTime         time.Time         `json:"end_time"`
	DurationMs      int64             `json:"duration_ms"`
	NumBidders      int               `json:"num_bidders_called"`
	NumBidsReceived int               `json:"num_bids_received"`
	Winner          Bid               `json:"winner,omitempty"`
	Bids            []Bid             `json:"bids"`
	Attributes      map[string]string `json:"attributes"`
}

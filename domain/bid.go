package domain

// Bid represents a bidder response
type Bid struct {
	BidderID  int     `json:"bidder_id"`
	Price     float64 `json:"price"`
	LatencyMs int     `json:"latency_ms"`
}

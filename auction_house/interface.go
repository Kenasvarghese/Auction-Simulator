package auctionhouse

import "github.com/Kenasvarghese/Auction-Simulator/domain"

// AuctionHouse defines the interface for running auctions.
type AuctionHouse interface {
	RunAuction(auctionID int, timeoutMs int) domain.AuctionResult
}

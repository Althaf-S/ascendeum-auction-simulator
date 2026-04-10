package auction

import "time"

type Attribute map[string]int

type Auction struct {
	ID         int
	Attributes Attribute
	Timeout    time.Duration
}

type Bid struct {
	BidderID int
	Amount   float64
	Time     time.Time
}

type Result struct {
	AuctionID int
	WinnerID  int
	Amount    float64
	TotalBids int

	StartTime time.Time
	EndTime   time.Time

	MemoryBefore uint64
	MemoryAfter  uint64
}

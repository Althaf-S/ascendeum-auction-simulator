package auction

import (
	"context"
	"runtime"
	"time"
)

func RunAuction(a Auction) Result {

	startTime := time.Now()

	//Capture memory before runtime starts
	var memoryBefore runtime.MemStats
	runtime.ReadMemStats(&memoryBefore)

	ctx, cancel := context.WithTimeout(context.Background(), a.Timeout)
	defer cancel()

	//Buffered channel to handle burst bids (all send bids at same moment)
	bidch := make(chan Bid, 100)

	var bids []Bid

	// TODO : bid simulation

	for {
		select {
		case bid := <-bidch:
			bids = append(bids, bid)
		case <-ctx.Done():
			endTime := time.Now()

			var memoryAfter runtime.MemStats
			runtime.ReadMemStats(&memoryAfter)

			res := decideWinner(a.ID, bids)
			res.StartTime = startTime
			res.EndTime = endTime
			// Used .Alloc here becuase MemStats is a struct, and we need
			// a specific metric (current heap allocation) as uint64
			res.MemoryBefore = memoryBefore.Alloc
			res.MemoryAfter = memoryAfter.Alloc

			return res
		}
	}

}

func decideWinner(auctionID int, bids []Bid) Result {
	var maxBid float64

	// winnerID set as -1 because its a sentinal value as we set bidder ID as int for
	// this siumulation usecase. in case zero bidders and no winner exists
	// so no bid -1 will be bidder id.
	winnerId := -1

	for _, bid := range bids {
		if bid.Amount > maxBid {
			maxBid = bid.Amount
			winnerId = bid.BidderID
		}
	}

	return Result{
		AuctionID: auctionID,
		WinnerID:  winnerId,
		Amount:    maxBid,
		TotalBids: len(bids),
	}
}

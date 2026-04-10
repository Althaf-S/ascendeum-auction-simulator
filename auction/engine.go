package auction

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"ascendeum-auction-simulator/models"
)

func RunAuction(a models.Auction) models.Result {

	startTime := time.Now()

	//Capture memory before runtime starts
	var memoryBefore runtime.MemStats
	runtime.ReadMemStats(&memoryBefore)

	ctx, cancel := context.WithTimeout(context.Background(), a.Timeout)
	defer cancel()

	// Biding simulation

	numOfBidders := rand.Intn(100) + 1

	//Buffered channel to handle burst bids (some send bids at same moment)
	bidCh := make(chan models.Bid, numOfBidders) // range is 'numOfBidders' since simulation, not in real world  case

	var bids []models.Bid

	for i := 0; i < numOfBidders; i++ {
		bidderID := i + 1
		// Could introduce semaphore based workerpool when the resource constraints are too
		// tight to standardize resources much w.r.t vCPU
		go func(id int) {
			if rand.Float64() < 0.3 { //30% chance a bidder won't bid
				return
			}

			for { //simulates same bidder bidding multiple times
				// random delay between bids mimics network and other constraints, min 50ms
				time.Sleep(time.Millisecond * time.Duration(50+rand.Intn(200)))

				bid := models.Bid{
					BidderID: id,
					Amount:   float64(100 + rand.Intn(900)), // 100 to 999 amount
					Time:     time.Now(),
				}

				select { // safe send, will stop at timeout
				case bidCh <- bid:
				case <-ctx.Done():
					return
				}
			}
		}(bidderID)
	}

	for {
		select {
		case bid := <-bidCh:
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
			res.TotalBidders = numOfBidders

			return res
		}
	}

}

func decideWinner(auctionID int, bids []models.Bid) models.Result {
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

	return models.Result{
		AuctionID: auctionID,
		WinnerID:  winnerId,
		Amount:    maxBid,
		TotalBids: len(bids),
	}
}

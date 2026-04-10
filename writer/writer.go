package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ascendeum-auction-simulator/models"
)

// Ensure output directory exists
func ensureOutputDir() error {
	return os.MkdirAll("outputs", os.ModePerm)
}

// Write per-auction file
func WriteAuctionResult(res models.Result) error {
	if err := ensureOutputDir(); err != nil {
		return err
	}

	filename := filepath.Join("outputs", fmt.Sprintf("auction_%d.txt", res.AuctionID))

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	content := fmt.Sprintf(
		`Auction ID: %d
		Winner ID: %d
		Winning Amount: %.2f
		Total Bids: %d

		Start Time: %s
		End Time: %s
		Duration: %s

		Memory Before: %d bytes
		Memory After: %d bytes
		`,
		res.AuctionID,
		res.WinnerID,
		res.Amount,
		res.TotalBids,
		res.StartTime.Format(time.RFC3339),
		res.EndTime.Format(time.RFC3339),
		res.EndTime.Sub(res.StartTime),
		res.MemoryBefore,
		res.MemoryAfter,
	)

	_, err = file.WriteString(content)
	return err
}

// Write summary file

func WriteSummary(results []models.Result, totalTime time.Duration) error {
	if err := ensureOutputDir(); err != nil {
		return err
	}

	filename := filepath.Join("outputs", "summary.txt")

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var totalBids int
	var maxBid float64
	var minBid float64 = -1

	for _, r := range results {
		totalBids += r.TotalBids

		if r.Amount > maxBid {
			maxBid = r.Amount
		}

		if minBid == -1 || r.Amount < minBid {
			minBid = r.Amount
		}
	}

	avgBids := float64(totalBids) / float64(len(results))

	content := fmt.Sprintf(
		`Total Auctions: %d
		Total Bids: %d
		Average Bids per Auction: %.2f

		Highest Winning Bid: %.2f
		Lowest Winning Bid: %.2f

		Total Execution Time: %s
		`,
		len(results),
		totalBids,
		avgBids,
		maxBid,
		minBid,
		totalTime,
	)

	_, err = file.WriteString(content)
	return err
}

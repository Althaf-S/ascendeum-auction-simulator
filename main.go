package main

import (
	"fmt"
	"runtime"
	"time"

	"ascendeum-auction-simulator/auction"
	"ascendeum-auction-simulator/writer"
)

func main() {

	// Explicit CPU usage in current system without workerpool,could also be obtained from env
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Starting Auction Simulation...")

	start := time.Now()

	// Run all auctions
	results := auction.ManageAuctions()

	totalTime := time.Since(start)

	fmt.Println("Total time from start of the first auction and the completion of the last auction : ", totalTime)

	fmt.Println("Writing auction results to files, you can view the same on outputs directory...")

	// Write individual auction report and result
	for _, res := range results {
		err := writer.WriteAuctionResult(res)
		if err != nil {
			fmt.Println("Error writing auction file:", err)
		}
	}

	// Write summary of entire process
	err := writer.WriteSummary(results, totalTime)
	if err != nil {
		fmt.Println("Error writing summary:", err)
	}

	fmt.Println("Simulation complete.")
}

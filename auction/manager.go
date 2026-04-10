package auction

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func ManageAuctions() []Result {

	var wg sync.WaitGroup

	numOfAuction := 40
	resultsCh := make(chan Result, numOfAuction)

	startTime := time.Now()

	for i := 0; i < numOfAuction; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			auction := Auction{
				ID: id,
				Attributes: Attribute{
					"demand": rand.Intn(100),
					"value":  rand.Intn(100),
				},
				Timeout: time.Millisecond * time.Duration(500+rand.Intn(1000)),
			}

			result := RunAuction(auction)

			resultsCh <- result
		}(i + 1)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	results := make([]Result, 0, numOfAuction)
	for res := range resultsCh {
		results = append(results, res)
	}

	endTime := time.Now()

	totalTime := endTime.Sub(startTime)
	fmt.Println("Total execution time from start of first auction and end of last auction", totalTime)

	return results
}

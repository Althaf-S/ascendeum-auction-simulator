package auction

import (
	"math/rand"
	"sync"
	"time"

	"ascendeum-auction-simulator/models"
)

func ManageAuctions() []models.Result {

	var wg sync.WaitGroup

	numOfAuction := 40
	resultsCh := make(chan models.Result, numOfAuction)

	for i := 0; i < numOfAuction; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			auction := models.Auction{
				ID: id,
				Attributes: models.Attribute{
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

	results := make([]models.Result, 0, numOfAuction)
	for res := range resultsCh {
		results = append(results, res)
	}

	return results
}

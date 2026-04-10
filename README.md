# Ascendum Auction Simulator (Go)

> Concurrent auction system simulating bidders, timeouts, and resource-aware execution.

---

## Overview

This project implements a **high-performance auction simulator** in Go that:

* Runs **40 auctions concurrently**
* Simulates **up to 100 bidders per auction**
* Uses **timeout-driven auction closure**
* Sleep mechanism simulates network latency and bidder decision time based on auction attributes
* Determines winners based on **highest valid bid**
* Measures **execution time and memory usage**
* Writes results to files for analysis

This system is designed to balance **concurrency, correctness, and resource awareness** without over-engineering.

---

## Key Features

* **Concurrent auctions using goroutines**
* **Timeout based auction lifecycle**
* **Channel based bid collection**
* **Execution time & memory tracking**
* **Per-auction and summary file outputs**
* **Modular and clean architecture**

---

## Project Structure

```
ascendeum-auction-simulator/
├── go.mod
├── main.go
├── auction/
│   ├── engine.go
│   └── manager.go
├── models/
│   └── models.go
├── writer/
│   └── writer.go
```

---

## How It Works

### 1. Entry Point (`main.go`)

* Sets CPU usage:

```go
runtime.GOMAXPROCS(runtime.NumCPU())
```

* Starts simulation
* Writes results to files

---

### 2. Auction Manager (`manager.go`)

* Launches **40 auctions concurrently**
* Uses `sync.WaitGroup` for coordination and completion of same
* Aggregates results via channels

---

### 3. Auction Engine (`engine.go`)

Each auction:

1. Starts timer
2. Creates timeout context
3. Spawns bidders
4. Collects bids via channel
5. Stops on timeout
6. Determines winner
7. Records metrics

---

### 4. Bidder Behavior

* ~30% bidders skip participation
* Remaining bidders:

  * Place multiple bids
  * Send bids at random intervals
* Simulates real-world unpredictability

---

### 5. Communication via Channels

```go
bidCh := make(chan Bid, numOfBidders)
```

* Buffered channel prevents blocking
* Handles burst bid traffic efficiently

---

### 6. Timeout Handling

```go
ctx, cancel := context.WithTimeout(...)
```

* Ensures auction stops after deadline
* Prevents goroutine leaks

---

### 7. Safe Bid Submission

```go
select {
case bidCh <- bid:
case <-ctx.Done():
    return
}
```

* Guarantees no bids after timeout

---

### 8. Winner Selection

* Highest bid wins
* Simple and deterministic logic via comparison

---

### 9. Metrics Collected

* Execution time (per auction + total)
* Memory usage (`runtime.MemStats`)
* Total bids per auction

---

## Resource Standardization

### Objective

Ensure predictable behavior across machines with different CPU and RAM.

### Approach

* **CPU Awareness**

```go
runtime.GOMAXPROCS(runtime.NumCPU())
```

* **Controlled Execution**

  * Timeouts bound runtime
  * Buffered channels manage load

* **Design Note**

  * Worker pool can be added to strictly limit goroutines
  * Avoided here to preserve full concurrency and simplicity

---

## Output

### Per Auction

```
outputs/auction_<id>.txt
```

Includes:

* Winner details
* Total bids
* Execution duration
* Memory usage

---

### Summary File

```
outputs/summary.txt
```

Includes:

* Total auctions
* Total bids
* Average bids
* Highest & lowest winning bids
* Total execution time

---

## How to Run

```bash
# Install dependencies
go mod tidy

# Run the simulator
go run main.go
```

---

## Sample Output

```
Total Auctions: 40
Total Bids: 9444
Average Bids per Auction: 236.10

Highest Winning Bid: 999.00
Lowest Winning Bid: 966.00

Total Execution Time: ~1.45s
```

---

## Design Decisions

| Decision                | Reason                  |
| ----------------------- | ----------------------- |
| Goroutines for auctions | True parallel execution |
| Channels for bids       | Safe communication      |
| Context for timeout     | Clean cancellation      |
| Buffered channels       | Prevent blocking        |
| Random bidder logic     | Realistic simulation    |
| Separate models package | Maintainability         |

---

## Trade-offs

| Choice                | Trade-off                            |
| --------------------- | ------------------------------------ |
| No worker pool        | Simpler, but less strict CPU control |
| Random bidding        | Less deterministic results           |

---

## Future Improvements

* Add worker pool for stricter resource control
* Expose simulation via REST API
* Use JSON output instead of text
* Add structured logging
* Track CPU usage & goroutines

---

## Conclusion

This project demonstrates:

* Effective use of Go concurrency primitives
* Clean system architecture
* Resource-aware simulation design
* Real-world auction modeling

---

## Author

Built as part of a concurrency and system design exercise in Go.
